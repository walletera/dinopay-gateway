package payments

import (
    "context"
    "fmt"
    "net/http"
    "strconv"

    "github.com/google/uuid"
    dinopayApi "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/pkg/events/payments"
    paymentsApi "github.com/walletera/payments/api"
    "go.uber.org/zap"
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

    client, err := dinopayApi.NewClient("http://localhost:2090", dinopayApi.WithClient(http.DefaultClient))
    if err != nil {
        return fmt.Errorf("failed creating dinopay api client: %w", err)
    }

    dinopayResp, err := client.CreatePayment(context.Background(), &dinopayApi.Payment{
        Amount:   withdrawalCreated.Amount,
        Currency: withdrawalCreated.Currency,
        SourceAccount: dinopayApi.Account{
            AccountHolder: "hardcodedSourceAccountHolder",
            AccountNumber: "hardcodedSourceAccountNumber",
        },
        DestinationAccount: dinopayApi.Account{
            AccountHolder: withdrawalCreated.Beneficiary.Account.Holder,
            AccountNumber: strconv.Itoa(withdrawalCreated.Beneficiary.Account.Number),
        },
        CustomerTransactionId: dinopayApi.OptString{
            Value: withdrawalCreated.Id,
            Set:   true,
        },
    })

    if err != nil {
        return handleError(logger, "failed creating payment on dinopay", err)
    }

    if dinopayResp == nil {
        return handleError(logger, "dinopay response is nil", err)
    }

    payment, ok := dinopayResp.(*dinopayApi.Payment)
    if !ok {
        return handleError(logger, "unexpected dinopay response", err)
    }

    paymentsClient, err := paymentsApi.NewClient("http://localhost:2090")
    if err != nil {
        return handleError(logger, "failed creating payments api client", err)
    }

    withdrawalId, err := uuid.Parse(withdrawalCreated.Id)
    if err != nil {
        handleError(logger, "failed parsing WithdrawalCreated uuid", err)
    }

    _, err = paymentsClient.PatchWithdrawal(context.Background(),
        &paymentsApi.WithdrawalPatchBody{
            ExternalId: paymentsApi.OptUUID{
                Value: payment.ID.Value,
                Set:   true,
            },
            Status: paymentsApi.WithdrawalPatchBodyStatusConfirmed,
        },
        paymentsApi.PatchWithdrawalParams{
            WithdrawalId: withdrawalId,
        })

    if err != nil {
        handleError(logger, "failed updating withdrawal in payments service", err)
    }

    logger.Info("WithdrawalCreated event processed successfully", zap.String("withdrawal_id", withdrawalCreated.Id))

    return nil
}

func handleError(logger *zap.Logger, errMsg string, err error) error {
    logger.Error(errMsg, zap.Error(err))
    return fmt.Errorf("%s: %w", errMsg, err)
}
