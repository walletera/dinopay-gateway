package tests

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/cucumber/godog"
    "github.com/walletera/dinopay-gateway/internal/app"
    "github.com/walletera/dinopay-gateway/pkg/logs"
    "github.com/walletera/message-processor/pkg/events/payments"
    "github.com/walletera/message-processor/pkg/rabbitmq"
    msClient "github.com/walletera/mockserver-go-client/pkg/client"
    "go.uber.org/zap"
    "net/http"
    "net/url"
    "testing"
    "time"
)

const (
    rawWithdrawalCreatedEventKey                  = "rawWithdrawalCreatedEvent"
    logsWatcherKey                                = "logsWatcher"
    dinoPayEndpointCreatePaymentsExpectationIdKey = "dinoPayEndpointCreatePaymentsExpectationId"
    mockserverClientKey                           = "mockserverClient"
    loggerKey                                     = "logger"
)

type MockServerExpectation struct {
    ExpectationID string `json:"id"`
}

var logger, _ = zap.NewDevelopment()

func aRunningDinopayGateway(ctx context.Context) (context.Context, error) {
    go func() {
        err := app.NewApp().Run(ctx)
        if err != nil {
            logger.Error("failed running app", zap.Error(err))
        }
    }()

    return ctx, nil
}

func aWithdrawalCreatedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawWithdrawalCreatedEventKey, []byte(event.Content)), nil
}

func aDinoPayEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    if mockserverExpectation == nil || len(mockserverExpectation.Content) == 0 {
        return nil, fmt.Errorf("the mockserver expectation is empty or was not defined")
    }

    rawMockserverExpectation := []byte(mockserverExpectation.Content)

    mockserverUrl, err := url.Parse(fmt.Sprintf("http://localhost:%s", mockserverPort))
    if err != nil {
        return nil, fmt.Errorf("error building mockserver url: %w", err)
    }

    mockserverClient := msClient.NewClient(mockserverUrl, http.DefaultClient)

    ctx = context.WithValue(ctx, mockserverClientKey, mockserverClient)

    var unmarshalledExpectation MockServerExpectation
    err = json.Unmarshal(rawMockserverExpectation, &unmarshalledExpectation)
    if err != nil {
        fmt.Errorf("error unmarshalling expectation: %w", err)
    }

    ctx = context.WithValue(ctx, dinoPayEndpointCreatePaymentsExpectationIdKey, unmarshalledExpectation.ExpectationID)

    err = mockserverClient.CreateExpectation(ctx, rawMockserverExpectation)
    if err != nil {
        fmt.Errorf("error creating mockserver expectations")
    }

    return ctx, nil
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

func theDinopayGatewayProcessTheEventAndProduceTheFollowingLog(ctx context.Context, logMsg string) (context.Context, error) {
    logsWatcher := logsWatcherFromCtx(ctx)
    foundLogEntry := logsWatcher.WaitFor(logMsg, 10*time.Second)
    if !foundLogEntry {
        return ctx, fmt.Errorf("didn't find expected log entry")
    }
    return ctx, nil
}

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

func InitializeScenario(godogCtx *godog.ScenarioContext) {

    godogCtx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        logsWatcher := logs.NewWatcher()
        logsWatcher.Start()
        ctx = context.WithValue(ctx, logsWatcherKey, logsWatcher)
        logger, err := zap.NewDevelopment()
        if err != nil {
            return ctx, fmt.Errorf("failed creating zap logger: %w", err)
        }
        ctx = context.WithValue(ctx, loggerKey, logger)
        return ctx, nil
    })

    godogCtx.Given(`^a running dinopay-gateway$`, aRunningDinopayGateway)
    godogCtx.Given(`^a withdrawal created event:$`, aWithdrawalCreatedEvent)
    godogCtx.Given(`^a dinopay endpoint to create payments:$`, aDinoPayEndpointToCreatePayments)
    godogCtx.When(`^the event is published$`, theEventIsPublished)
    godogCtx.Then(`^the dinopay-gateway process the event and produce the following log:$`, theDinopayGatewayProcessTheEventAndProduceTheFollowingLog)

    godogCtx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

        verifyErr := verifyMockServerExpectations(ctx)
        if verifyErr != nil {
            return ctx, fmt.Errorf("the dinopay api was not called with the expected parameters: %w", verifyErr)
        }

        logsWatcher := logsWatcherFromCtx(ctx)
        logsWatcher.Stop()

        mockserverClient := mockserverClientFromCtx(ctx)
        clearReqErr := mockserverClient.Clear(ctx)
        if clearReqErr != nil {
            return nil, fmt.Errorf("failed clearing mockserver: %w", clearReqErr)
        }
        return ctx, nil
    })
}

func logsWatcherFromCtx(ctx context.Context) *logs.Watcher {
    return ctx.Value(logsWatcherKey).(*logs.Watcher)
}

func mockserverClientFromCtx(ctx context.Context) *msClient.Client {
    return ctx.Value(mockserverClientKey).(*msClient.Client)
}

func expectationIdFromCtx(ctx context.Context) string {
    return ctx.Value(dinoPayEndpointCreatePaymentsExpectationIdKey).(string)
}

func verifyMockServerExpectations(ctx context.Context) error {
    mockserverClient := mockserverClientFromCtx(ctx)
    verificationErr := mockserverClient.VerifyRequest(ctx, msClient.VerifyRequestBody{
        ExpectationId: msClient.ExpectationId{
            Id: expectationIdFromCtx(ctx),
        },
    })
    if verificationErr != nil {
        return verificationErr
    }
    return nil
}
