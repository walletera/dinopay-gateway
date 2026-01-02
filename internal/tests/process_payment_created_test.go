package tests

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/cucumber/godog"
    "github.com/walletera/dinopay-gateway/internal/app"
    "github.com/walletera/eventskit/events"

    "github.com/walletera/eventskit/rabbitmq"
)

const (
    rawWithdrawalCreatedEventKey                     = "rawWithdrawalCreatedEvent"
    dinoPayEndpointCreatePaymentsExpectationIdKey    = "dinoPayEndpointCreatePaymentsExpectationId"
    paymentsEndpointUpdateWithdrawalExpectationIdKey = "paymentsEndpointUpdateWithdrawalExpectationId"
    expectationTimeout                               = 5 * time.Second
)

func TestPaymentCreatedEventProcessing(t *testing.T) {

    suite := godog.TestSuite{
        ScenarioInitializer: InitializeProcessWithdrawalCreatedScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features/payment_created.feature"},
            TestingT: t, // Testing instance that will run subtests.
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeProcessWithdrawalCreatedScenario(ctx *godog.ScenarioContext) {
    ctx.Before(beforeScenarioHook)
    ctx.Given(`^a running dinopay-gateway$`, aRunningDinopayGateway)
    ctx.Given(`^a PaymentCreated event:$`, aPaymentCreatedEvent)
    ctx.Given(`^a dinopay endpoint to create payments:$`, aDinopayEndpointToCreatePayments)
    ctx.Given(`^a payments endpoint to update payments:$`, aPaymentsEndpointToUpdatePayments)
    ctx.When(`^the event is published$`, theEventIsPublished)
    ctx.Then(`^the dinopay-gateway creates the corresponding payment on the DinoPay API$`, theDinopayGatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI)
    ctx.Then(`^the dinopay-gateway updates the payment on payments service$`, theDinopayGatewayUpdatesThePaymentOnPaymentsService)
    ctx.Then(`the dinopay-gateway fails creating the corresponding payment on the DinoPay API$`, theDinoPayGatewayFailsCreatingTheCorrespondingPayment)
    ctx.Then(`^the dinopay-gateway produces the following log:$`, theDinopayGatewayProducesTheFollowingLog)
    ctx.After(afterScenarioHook)
}

func aPaymentCreatedEvent(ctx context.Context, eventFilePath *godog.DocString) (context.Context, error) {
    return context.WithValue(ctx, rawWithdrawalCreatedEventKey, readFile(eventFilePath)), nil
}

func aDinopayEndpointToCreatePayments(ctx context.Context, mockserverExpectationFilePath *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectationFilePath, dinoPayEndpointCreatePaymentsExpectationIdKey)
}

func aPaymentsEndpointToUpdatePayments(ctx context.Context, mockserverExpectationFilePath *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectationFilePath, paymentsEndpointUpdateWithdrawalExpectationIdKey)
}

func theEventIsPublished(ctx context.Context) (context.Context, error) {
    publisher, err := rabbitmq.NewClient(
        rabbitmq.WithExchangeName(app.RabbitMQPaymentsExchangeName),
        rabbitmq.WithExchangeType(app.RabbitMQExchangeType),
    )
    if err != nil {
        return nil, fmt.Errorf("error creating rabbitmq client: %s", err.Error())
    }

    rawEvent := ctx.Value(rawWithdrawalCreatedEventKey).([]byte)
    err = publisher.Publish(ctx, publishable{rawEvent: rawEvent}, events.RoutingInfo{
        Topic:      app.RabbitMQPaymentsExchangeName,
        RoutingKey: app.RabbitMQPaymentCreatedRoutingKey,
    })
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

func theDinopayGatewayUpdatesThePaymentOnPaymentsService(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, paymentsEndpointUpdateWithdrawalExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinoPayGatewayFailsCreatingTheCorrespondingPayment(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, dinoPayEndpointCreatePaymentsExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}
