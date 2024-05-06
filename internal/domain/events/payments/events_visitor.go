package payments

import (
    "context"
    "fmt"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events"
    dinopayEvents "github.com/walletera/dinopay-gateway/internal/domain/events/dinopay"
    "github.com/walletera/dinopay-gateway/internal/domain/ports/output/dinopay"
    dinopayApi "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/payments"
    "go.uber.org/zap"
)

type EventsVisitor struct {
    dinopayClient dinopay.Client
    esDB          events.DB
}

func NewEventsVisitor(
    dinopayClient dinopay.Client,
    esDB events.DB,
) *EventsVisitor {
    return &EventsVisitor{
        dinopayClient: dinopayClient,
        esDB:          esDB,
    }
}

func (e *EventsVisitor) VisitWithdrawalCreated(withdrawalCreated payments.WithdrawalCreatedEvent) error {
    logger, err := zap.NewDevelopment()
    if err != nil {
        return fmt.Errorf("failed creating logger: %w", err)
    }

    dinopayResp, err := e.dinopayClient.CreatePayment(context.Background(), &dinopayApi.Payment{
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

    withdrawalId, err := uuid.Parse(withdrawalCreated.Id)
    if err != nil {
        handleError(logger, "failed parsing WithdrawalCreated uuid", err)
    }

    outboundPaymentCreated := dinopayEvents.OutboundPaymentCreated{
        Id:                   uuid.New(),
        WithdrawalId:         withdrawalId,
        DinopayPaymentId:     payment.ID.Value,
        DinopayPaymentStatus: string(payment.Status.Value),
    }

    err = e.esDB.AppendEvents(dinopayEvents.BuildStreamName(payment.ID.Value.String()), outboundPaymentCreated)
    if err != nil {
        return fmt.Errorf("failed adding OutboundPaymentCreated event to the repository: %w", err)
    }

    logger.Info("WithdrawalCreated event processed successfully", zap.String("withdrawal_id", withdrawalCreated.Id))

    return nil
}

func handleError(logger *zap.Logger, errMsg string, err error) error {
    logger.Error(errMsg, zap.Error(err))
    return fmt.Errorf("%s: %w", errMsg, err)
}
