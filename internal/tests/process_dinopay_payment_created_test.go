package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/walletera/dinopay-gateway/internal/app"
)

const (
	rawDinopayPaymentCreatedEventKey            = "rawDinopayPaymentCreatedEventKey"
	accountsGetAccountEndpointExpectationKey    = "accountsGetAccountEndpointExpectationKey"
	paymentsCreateDepositEndpointExpectationKey = "paymentsCreateDepositEndpointExpectationKey"
)

func TestDinopayPaymentCreatedEventProcessing(t *testing.T) {

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeProcessDinopayPaymentCreatedScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/dinopay_payment_created.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeProcessDinopayPaymentCreatedScenario(ctx *godog.ScenarioContext) {
	ctx.Before(beforeScenarioHook)
	ctx.Step(`^a running dinopay-gateway$`, aRunningDinopayGateway)
	ctx.Step(`^a DinoPay PaymentCreated event:$`, aDinoPayPaymentCreatedEvent)
	ctx.Step(`^an accounts endpoint to get accounts:$`, anAccountsEndpointToGetAccounts)
	ctx.Step(`^a payments endpoint to create payments:$`, aPaymentsEndpointToCreatePayments)
	ctx.When(`^the webhook event is received$`, theWebhookEventIsReceived)
	ctx.Step(`^the dinopay-gateway creates the corresponding payment on the Payments API$`, theDinopaygatewayCreatesTheCorrespondingPaymentOnThePaymentsAPI)
	ctx.Step(`^the dinopay-gateway produces the following log:$`, theDinopayGatewayProducesTheFollowingLog)
	ctx.After(afterScenarioHook)
}

<<<<<<< Updated upstream
func aDinoPayPaymentCreatedEvent(ctx context.Context, jsonEventFilePath *godog.DocString) (context.Context, error) {
    return context.WithValue(ctx, rawDinopayPaymentCreatedEventKey, readFile(jsonEventFilePath)), nil
}

func aPaymentsEndpointToCreateDeposits(ctx context.Context, mockserverExpectationFilePath *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectationFilePath, paymentsCreateDepositEndpointExpectationKey)
=======
func aDinoPayPaymentCreatedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
	if event == nil || len(event.Content) == 0 {
		return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
	}
	return context.WithValue(ctx, rawDinopayPaymentCreatedEventKey, []byte(event.Content)), nil
}

func anAccountsEndpointToGetAccounts(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
	return createMockServerExpectation(ctx, mockserverExpectation, accountsGetAccountEndpointExpectationKey)
}

func aPaymentsEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
	return createMockServerExpectation(ctx, mockserverExpectation, paymentsCreateDepositEndpointExpectationKey)
>>>>>>> Stashed changes
}

func theWebhookEventIsReceived(ctx context.Context) (context.Context, error) {
	rawEvent := ctx.Value(rawDinopayPaymentCreatedEventKey).([]byte)
	url := fmt.Sprintf("http://127.0.0.1:%d/webhooks", app.WebhookServerPort)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rawEvent))
	if err != nil {
		return ctx, fmt.Errorf("failed sending webhook event: %w", err)
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return ctx, fmt.Errorf("failed sending request to payments api: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return ctx, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
	return ctx, nil
}

func theDinopaygatewayCreatesTheCorrespondingPaymentOnThePaymentsAPI(ctx context.Context) (context.Context, error) {
	id := expectationIdFromCtx(ctx, paymentsCreateDepositEndpointExpectationKey)
	err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
	return ctx, err
}
