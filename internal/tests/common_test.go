package tests

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "net/url"
    "time"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/cucumber/godog"
    "github.com/walletera/dinopay-gateway/internal/app"
    "github.com/walletera/eventskit/eventstoredb"
    "github.com/walletera/eventskit/rabbitmq"
    slogwatcher "github.com/walletera/logs-watcher/slog"
    msClient "github.com/walletera/mockserver-go-client/pkg/client"
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
    "go.uber.org/zap/zapcore"
    "golang.org/x/sync/errgroup"
)

const (
    mockserverUrl             = "http://localhost:2090"
    eventStoreDBUrl           = "esdb://localhost:2113?tls=false"
    appKey                    = "app"
    appCtxCancelFuncKey       = "appCtxCancelFuncKey"
    logsWatcherKey            = "logsWatcher"
    logsWatcherWaitForTimeout = 5 * time.Second
)

type MockServerExpectation struct {
    ExpectationID string `json:"id"`
}

func beforeScenarioHook(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
    handler, err := newZapHandler()
    if err != nil {
        return ctx, err
    }
    logsWatcher := slogwatcher.NewWatcher(handler)
    ctx = context.WithValue(ctx, logsWatcherKey, logsWatcher)
    return ctx, nil
}

func afterScenarioHook(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {

    clearReqErr := mockServerClient().Clear(ctx)
    if clearReqErr != nil {
        return nil, fmt.Errorf("failed clearing mockserver: %w", clearReqErr)
    }

    logsWatcher := logsWatcherFromCtx(ctx)

    appCtxCancelFuncFromCtx(ctx)()
    appFromCtx(ctx).Stop(ctx)
    foundLogEntry := logsWatcher.WaitFor("dinopay-gateway stopped", logsWatcherWaitForTimeout)
    if !foundLogEntry {
        return ctx, fmt.Errorf("app termination failed (didn't find expected log entry)")
    }

    err = logsWatcher.Stop()
    if err != nil {
        return ctx, fmt.Errorf("failed stopping the logsWatcher: %w", err)
    }

    return ctx, nil
}

func aRunningDinopayGateway(ctx context.Context) (context.Context, error) {

    ctx, err := esdbByCategoryProjectionEnabled(ctx)
    if err != nil {
        return ctx, fmt.Errorf("failed enabling by-category projection: %w", err)
    }

    ctx, err = anEventstoreDBPersistentSubscriptionForCategory(ctx, app.ESDB_ByCategoryProjection_OutboundPayment)
    if err != nil {
        return ctx, fmt.Errorf("failed creating persistent subscription for $ce-outboundPayment category on esdb: %w", err)
    }

    ctx, err = anEventstoreDBPersistentSubscriptionForCategory(ctx, app.ESDB_ByCategoryProjection_InboundPayment)
    if err != nil {
        return ctx, fmt.Errorf("failed creating persistent subscription for $ce-inboundPayment category on esdb: %w", err)
    }

    logHandler := logsWatcherFromCtx(ctx).DecoratedHandler()

    appCtx, appCtxCancelFunc := context.WithCancel(ctx)
    app, err := app.NewApp(
        app.WithRabbitmqHost(rabbitmq.DefaultHost),
        app.WithRabbitmqPort(rabbitmq.DefaultPort),
        app.WithRabbitmqUser(rabbitmq.DefaultUser),
        app.WithRabbitmqPassword(rabbitmq.DefaultPassword),
        app.WithDinopayUrl(mockserverUrl),
        app.WithPaymentsUrl(mockserverUrl),
        app.WithESDBUrl(eventStoreDBUrl),
        app.WithLogHandler(logHandler),
    )
    if err != nil {
        panic("failed initializing app: " + err.Error())
    }

    err = app.Run(appCtx)
    if err != nil {
        panic("failed running app" + err.Error())
    }

    ctx = context.WithValue(ctx, appKey, app)
    ctx = context.WithValue(ctx, appCtxCancelFuncKey, appCtxCancelFunc)

    foundLogEntry := logsWatcherFromCtx(ctx).WaitFor("dinopay-gateway started", logsWatcherWaitForTimeout)
    if !foundLogEntry {
        return ctx, fmt.Errorf("app startup failed (didn't find expected log entry)")
    }

    return ctx, nil
}

func esdbByCategoryProjectionEnabled(ctx context.Context) (context.Context, error) {
    err := eventstoredb.EnableByCategoryProjection(ctx, eventStoreDBUrl)
    if err != nil {
        return ctx, err
    }
    return ctx, nil
}

func anEventstoreDBPersistentSubscriptionForCategory(ctx context.Context, categoryName string) (context.Context, error) {
    subscriptionSettings := esdb.SubscriptionSettingsDefault()
    subscriptionSettings.ResolveLinkTos = true

    err := eventstoredb.CreatePersistentSubscription(
        eventStoreDBUrl,
        categoryName,
        app.ESDB_SubscriptionGroupName,
        subscriptionSettings,
    )
    if err != nil {
        return ctx, err
    }

    return ctx, nil
}

func theDinopayGatewayProducesTheFollowingLog(ctx context.Context, logMsg string) (context.Context, error) {
    logsWatcher := logsWatcherFromCtx(ctx)
    foundLogEntry := logsWatcher.WaitFor(logMsg, logsWatcherWaitForTimeout)
    if !foundLogEntry {
        return ctx, fmt.Errorf("didn't find expected log entry")
    }
    return ctx, nil
}

func createMockServerExpectation(ctx context.Context, mockserverExpectation *godog.DocString, ctxKey string) (context.Context, error) {
    if mockserverExpectation == nil || len(mockserverExpectation.Content) == 0 {
        return nil, fmt.Errorf("the mockserver expectation is empty or was not defined")
    }

    rawMockserverExpectation := []byte(mockserverExpectation.Content)

    var unmarshalledExpectation MockServerExpectation
    err := json.Unmarshal(rawMockserverExpectation, &unmarshalledExpectation)
    if err != nil {
        fmt.Errorf("error unmarshalling expectation: %w", err)
    }

    ctx = context.WithValue(ctx, ctxKey, unmarshalledExpectation.ExpectationID)

    err = mockServerClient().CreateExpectation(ctx, rawMockserverExpectation)
    if err != nil {
        fmt.Errorf("error creating mockserver expectations")
    }

    return ctx, nil
}

func mockServerClient() *msClient.Client {
    mockserverUrl, err := url.Parse(fmt.Sprintf("http://localhost:%s", mockserverPort))
    if err != nil {
        panic("error building mockserver url: " + err.Error())
    }

    return msClient.NewClient(mockserverUrl, http.DefaultClient)
}

func appCtxCancelFuncFromCtx(ctx context.Context) context.CancelFunc {
    return ctx.Value(appCtxCancelFuncKey).(context.CancelFunc)
}

func expectationIdFromCtx(ctx context.Context, ctxKey string) string {
    return ctx.Value(ctxKey).(string)
}

func logsWatcherFromCtx(ctx context.Context) *slogwatcher.Watcher {
    return ctx.Value(logsWatcherKey).(*slogwatcher.Watcher)
}

func appFromCtx(ctx context.Context) *app.App {
    return ctx.Value(appKey).(*app.App)
}

func verifyExpectationMetWithin(ctx context.Context, expectationID string, timeout time.Duration) error {
    errGroup := new(errgroup.Group)
    timeoutCh := time.After(timeout)
    errGroup.Go(func() error {
        var err error
        for {
            select {
            case <-timeoutCh:
                return fmt.Errorf("expectation %s was not met whithin %s: %w", expectationID, timeout.String(), err)
            default:
                err = verifyExpectationMet(ctx, expectationID)
                if err == nil {
                    return nil
                }
                time.Sleep(1 * time.Second)
            }
        }
    })
    return errGroup.Wait()
}

func verifyExpectationMet(ctx context.Context, expectationID string) error {
    verificationErr := mockServerClient().VerifyRequest(ctx, msClient.VerifyRequestBody{
        ExpectationId: msClient.ExpectationId{
            Id: expectationID,
        },
    })
    if verificationErr != nil {
        return verificationErr
    }
    return nil
}

func newZapHandler() (slog.Handler, error) {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
    zapConfig := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:    encoderConfig,
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
    zapLogger, err := zapConfig.Build()
    if err != nil {
        return nil, err
    }
    return zapslog.NewHandler(
        zapLogger.Core(),
        // never add stacktrace
        zapslog.AddStacktraceAt(slog.LevelError+1),
        ), nil
}
