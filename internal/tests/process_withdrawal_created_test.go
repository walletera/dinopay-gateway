package tests

import (
    "testing"

    "github.com/cucumber/godog"
)

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
    ctx.Given(`^a withdrawal created event:$`, aWithdrawalCreatedEvent)
    ctx.Given(`^a dinopay endpoint to create payments:$`, aDinopayEndpointToCreatePayments)
    ctx.When(`^the event is published$`, theEventIsPublished)
    ctx.Then(`^the dinopay-gateway creates the corresponding payment on the DinoPay API$`, theDinopaygatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI)
    ctx.Then(`^the dinopay-gateway produces the following log:$`, theDinopaygatewayProducesTheFollowingLog)
}

func aWithdrawalCreatedEvent(arg1 *godog.DocString) error {
    return godog.ErrPending
}

func aDinopayEndpointToCreatePayments(arg1 *godog.DocString) error {
    return godog.ErrPending
}

func theEventIsPublished() error {
    return godog.ErrPending
}

func theDinopaygatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI() error {
    return godog.ErrPending
}

func theDinopaygatewayProducesTheFollowingLog(arg1 *godog.DocString) error {
    return godog.ErrPending
}
