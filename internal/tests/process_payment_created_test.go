package tests

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/cucumber/godog"
    "github.com/walletera/message-processor/payments"
    "github.com/walletera/message-processor/rabbitmq"
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

func aPaymentCreatedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawWithdrawalCreatedEventKey, []byte(event.Content)), nil
}

func aDinopayEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, dinoPayEndpointCreatePaymentsExpectationIdKey)
}

func aPaymentsEndpointToUpdatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
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
