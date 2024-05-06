package tests

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "testing"
    "time"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/cucumber/godog"
    "github.com/walletera/dinopay-gateway/internal/app"
    "github.com/walletera/dinopay-gateway/pkg/eventstoredb"
    "github.com/walletera/logs-watcher"
    "github.com/walletera/message-processor/payments"
    "github.com/walletera/message-processor/rabbitmq"
    msClient "github.com/walletera/mockserver-go-client/pkg/client"
    "go.uber.org/zap"
    "golang.org/x/sync/errgroup"
)

const (
    mockserverUrl                                    = "http://localhost:2090"
    eventStoreDBUrl                                  = "esdb://localhost:2113?tls=false"
    appCtxCancelFuncKey                              = "appCtxCancelFuncKey"
    rawWithdrawalCreatedEventKey                     = "rawWithdrawalCreatedEvent"
    dinoPayEndpointCreatePaymentsExpectationIdKey    = "dinoPayEndpointCreatePaymentsExpectationId"
    paymentsEndpointUpdateWithdrawalExpectationIdKey = "paymentsEndpointUpdateWithdrawalExpectationId"
    expectationTimeout                               = 5 * time.Second
    logsWatcherKey                                   = "logsWatcher"
)

type MockServerExpectation struct {
    ExpectationID string `json:"id"`
}

var logger, _ = zap.NewDevelopment()

func TestFeatures(t *testing.T) {

    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features"},
            TestingT: t, // Testing instance that will run subtests.
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeScenario(ctx *godog.ScenarioContext) {

    ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        logsWatcher := logs.NewWatcher()
        logsWatcher.Start()
        ctx = context.WithValue(ctx, logsWatcherKey, logsWatcher)
        return ctx, nil
    })

    ctx.Given(`^a running dinopay-gateway$`, aRunningDinopayGateway)
    ctx.Given(`^a withdrawal created event:$`, aWithdrawalCreatedEvent)
    ctx.Given(`^a dinopay endpoint to create payments:$`, aDinopayEndpointToCreatePayments)
    ctx.Given(`^a payments endpoint to update withdrawals:$`, aPaymentsEndpointToUpdateWithdrawals)
    ctx.When(`^the event is published$`, theEventIsPublished)
    ctx.Then(`^the dinopay-gateway creates the corresponding payment on the DinoPay API$`, theDinopayGatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI)
    ctx.Then(`^the dinopay-gateway updates the withdrawal on payments service$`, theDinopayGatewayUpdatesTheWithdrawalOnPaymentsService)
    ctx.Then(`the dinopay-gateway fails creating the corresponding payment on the DinoPay API$`, theDinoPayGatewayFailsCreatingTheCorrespondingPayment)
    ctx.Then(`^the dinopay-gateway produces the following log:$`, theDinopayGatewayProducesTheFollowingLog)

    ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

        clearReqErr := mockServerClient().Clear(ctx)
        if clearReqErr != nil {
            return nil, fmt.Errorf("failed clearing mockserver: %w", clearReqErr)
        }

        logsWatcher := logsWatcherFromCtx(ctx)

        appCtxCancelFuncFromCtx(ctx)()
        foundLogEntry := logsWatcher.WaitFor("dinopay-gateway stopped", 5*time.Second)
        if !foundLogEntry {
            return ctx, fmt.Errorf("didn't find expected log entry")
        }

        err = logsWatcher.Stop()
        if err != nil {
            return ctx, fmt.Errorf("failed stopping the logsWatcher: %w", err)
        }

        return ctx, nil
    })
}

func aRunningDinopayGateway(ctx context.Context) (context.Context, error) {

    ctx, err := esdbByCategoryProjectionEnabled(ctx)
    if err != nil {
        return ctx, fmt.Errorf("failed enabling by-category projection: %w", err)
    }

    ctx, err = anEventstoreDBPersistentSubscriptionForCategoryOutboundPayment(ctx)
    if err != nil {
        return ctx, fmt.Errorf("failed creating persistent subscription on esdb: %w", err)
    }

    appCtx, appCtxCancelFunc := context.WithCancel(ctx)
    go func() {
        err := app.NewApp(
            app.WithDinopayUrl(mockserverUrl),
            app.WithPaymentsUrl(mockserverUrl),
            app.WithESDBUrl(eventStoreDBUrl),
        ).Run(appCtx)
        if err != nil {
            logger.Error("failed running app", zap.Error(err))
        }
    }()

    ctx = context.WithValue(ctx, appCtxCancelFuncKey, appCtxCancelFunc)

    foundLogEntry := logsWatcherFromCtx(ctx).WaitFor("dinopay-gateway started", 5*time.Second)
    if !foundLogEntry {
        return ctx, fmt.Errorf("didn't find expected log entry")
    }

    return ctx, nil
}

func esdbByCategoryProjectionEnabled(ctx context.Context) (context.Context, error) {
    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        fmt.Sprintf("http://127.0.0.1:%s/projection/$by_category/command/enable", eventStoreDBPort),
        nil,
    )
    if err != nil {
        return ctx, fmt.Errorf("failed creating request for enabling $by_category projection: %w", err)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Length", "0")
    _, err = http.DefaultClient.Do(req)
    if err != nil {
        return ctx, fmt.Errorf("failed enabling $by_category projection: %w", err)
    }
    return ctx, nil
}

func anEventstoreDBPersistentSubscriptionForCategoryOutboundPayment(ctx context.Context) (context.Context, error) {
    subscriptionSettings := esdb.SubscriptionSettingsDefault()
    subscriptionSettings.ResolveLinkTos = true
    subscriptionSettings.MaxRetryCount = 3

    esdbClient, err := eventstoredb.GetESDBClient(eventStoreDBUrl)
    if err != nil {
        return ctx, err
    }

    err = esdbClient.CreatePersistentSubscription(
        context.Background(),
        app.ESDB_ByCategoryProjection_OutboundPayment,
        app.ESDB_SubscriptionGroupName,
        esdb.PersistentStreamSubscriptionOptions{
            Settings: &subscriptionSettings,
        },
    )
    // FIXME: delete persistent subscription on the After hook
    if err != nil {
        var esdbError *esdb.Error
        ok := errors.As(err, &esdbError)
        if !ok || !esdbError.IsErrorCode(esdb.ErrorCodeResourceAlreadyExists) {
            return ctx, fmt.Errorf("CreatePersistentSubscription failed: %w", err)
        }
    }

    return ctx, nil
}

func aWithdrawalCreatedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawWithdrawalCreatedEventKey, []byte(event.Content)), nil
}

func aDinopayEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, dinoPayEndpointCreatePaymentsExpectationIdKey)
}

func aPaymentsEndpointToUpdateWithdrawals(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, paymentsEndpointUpdateWithdrawalExpectationIdKey)
}

func theEventIsPublished(ctx context.Context) (context.Context, error) {
    publisher, err := rabbitmq.NewClient(
        rabbitmq.WithExchangeName(payments.RabbitMQExchangeName),
        rabbitmq.WithExchangeType(payments.RabbitMQExchangeType),
    )
    if err != nil {
        return nil, fmt.Errorf("error creating rabbitmq client: %s", err.Error())
    }

    rawEvent := ctx.Value(rawWithdrawalCreatedEventKey).([]byte)
    err = publisher.Publish(ctx, rawEvent, payments.RabbitMQRoutingKey)
    if err != nil {
        return nil, fmt.Errorf("error publishing WithdrawalCreated event to rabbitmq: %s", err.Error())
    }

    return ctx, nil
}

func theDinopayGatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, dinoPayEndpointCreatePaymentsExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinopayGatewayUpdatesTheWithdrawalOnPaymentsService(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, paymentsEndpointUpdateWithdrawalExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinoPayGatewayFailsCreatingTheCorrespondingPayment(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, dinoPayEndpointCreatePaymentsExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinopayGatewayProducesTheFollowingLog(ctx context.Context, logMsg string) (context.Context, error) {
    logsWatcher := logsWatcherFromCtx(ctx)
    foundLogEntry := logsWatcher.WaitFor(logMsg, 10*time.Second)
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

func logsWatcherFromCtx(ctx context.Context) *logs.Watcher {
    return ctx.Value(logsWatcherKey).(*logs.Watcher)
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
