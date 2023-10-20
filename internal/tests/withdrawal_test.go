package tests

import (
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/walletera/dinopay-gateway/internal/app"
    "github.com/walletera/dinopay-gateway/internal/tests/expectations/dinopay"
    "github.com/walletera/dinopay-gateway/internal/tests/testdata"
    "github.com/walletera/message-processor/pkg/events/payments"
    "github.com/walletera/message-processor/pkg/rabbitmq"
    "testing"
    "time"
)

/*
   Given a WithdrawalCreated event published by payments service
    When dinopay-gateway process the event
    Then a Payment is created on DinoPay with the Beneficiary Account as the Destination Account
     And the Withdrawal status is updated in payments service to the `delivered` status
*/

func Test_ProcessWithdrawalCreatedEvent(t *testing.T) {
    integrationTest := newIntegrationTest(
        t,
        "process-withdrawal-created-event",
        dinopay.CreatePaymentSucceedExpectation,
    )

    integrationTest.Run(func(ctx context.Context, t *testing.T) {

        logsWatcher := NewLogsWatcher()
        logsWatcher.Start()

        go func() {
            app := app.NewApp()
            err := app.Run(ctx)
            if err != nil {
                t.Errorf("failed running app: %s", err.Error())
            }
        }()

        foundLogEntry := logsWatcher.WaitFor("dinopay-gateway started", 5*time.Second)
        if !foundLogEntry {
            t.Error("didn't find expected log entry")
            return
        }

        publisher, err := rabbitmq.NewClient(
            rabbitmq.WithExchangeName(payments.RabbitMQExchangeName),
            rabbitmq.WithExchangeType(payments.RabbitMQExchangeType),
        )
        if err != nil {
            t.Errorf("error creating rabbitmq client: %s", err.Error())
            return
        }

        err = publisher.Publish(ctx, testdata.WithdrawalCreated, payments.RabbitMQRoutingKey)
        if err != nil {
            t.Errorf("error publishing WithdrawalCreated event to rabbitmq: %s", err.Error())
            return
        }

        foundLogEntry = logsWatcher.WaitFor("processing WithdrawalCreated event", 5*time.Second)
        assert.True(t, foundLogEntry, "didn't find expected log entry")
    })
}
