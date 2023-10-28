package payments

import (
    "context"
    "fmt"
    "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/pkg/events/payments"
    "go.uber.org/zap"
    "net/http"
    "strconv"
)

type EventsVisitor struct {
}

func NewEventsVisitor() *EventsVisitor {
    return &EventsVisitor{}
}

func (e *EventsVisitor) VisitWithdrawalCreated(withdrawalCreated payments.WithdrawalCreated) error {
    logger, err := zap.NewDevelopment()
    if err != nil {
        return fmt.Errorf("failed creating logger: %w", err)
    }

    client, err := api.NewClient("http://localhost:2090", api.WithClient(http.DefaultClient))
    if err != nil {
        return fmt.Errorf("failed creating dinopay api client: %w", err)
    }

    _, err = client.CreatePayment(context.Background(), &api.Payment{
        Amount:   withdrawalCreated.Amount,
        Currency: withdrawalCreated.Currency,
        SourceAccount: api.Account{
            AccountHolder: "hardcodedSourceAccountHolder",
            AccountNumber: "hardcodedSourceAccountNumber",
        },
        DestinationAccount: api.Account{
            AccountHolder: withdrawalCreated.Beneficiary.Account.Holder,
            AccountNumber: strconv.Itoa(withdrawalCreated.Beneficiary.Account.Number),
        },
        CustomerTransactionId: api.OptString{
            Value: withdrawalCreated.Id,
            Set:   true,
        },
    })
    if err != nil {
        return fmt.Errorf("failed creating payment on dinopay: %w", err)
    }

    logger.Info("WithdrawalCreated event processed successfully", zap.String("withdrawal_id", withdrawalCreated.Id))

    return nil
}
