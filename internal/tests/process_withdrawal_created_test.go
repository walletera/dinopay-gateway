package tests

import (
    "github.com/cucumber/godog"
    "testing"
)

func aWithdrawalCreatedEvent(event *godog.DocString) error {
    return godog.ErrPending
}

func theEventIsPublished() error {
    return godog.ErrPending
}

func dinoPayGatewayProcessTheEvent() error {
    return godog.ErrPending
}

func theDinoPayApiIsCalledWithTheCorrectParameters(expectation *godog.DocString) error {
    return godog.ErrPending
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

func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Given(`^a withdrawal created event:$`, aWithdrawalCreatedEvent)
    ctx.When(`^the event is published$`, theEventIsPublished)
    ctx.Then(`^dinopay-gateway process the event$`, dinoPayGatewayProcessTheEvent)
    ctx.Then(`^dinopay api is called with the correct parameters:$`, theDinoPayApiIsCalledWithTheCorrectParameters)
}
