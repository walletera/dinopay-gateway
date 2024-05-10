package payments

import (
    "context"
    "fmt"
    "log/slog"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events"
    dinopayEvents "github.com/walletera/dinopay-gateway/internal/domain/events/dinopay"
    "github.com/walletera/dinopay-gateway/internal/domain/ports/output/dinopay"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    dinopayApi "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/payments"
)

type EventsVisitor struct {
    dinopayClient dinopay.Client
    esDB          events.DB
    logger        *slog.Logger
}

func NewEventsVisitor(dinopayClient dinopay.Client, esDB events.DB, logger *slog.Logger, ) *EventsVisitor {
    return &EventsVisitor{
        dinopayClient: dinopayClient,
        esDB:          esDB,
        logger:        logger.With(logattr.Component("payments/EventsVisitor")),
    }
}

func (ev *EventsVisitor) VisitWithdrawalCreated(withdrawalCreated payments.WithdrawalCreatedEvent) error {
    logger := ev.logger.With(
        logattr.EventType("WithdrawalCreated"),
        logattr.WithdrawalId(withdrawalCreated.Id),
    )

    // TODO modify visitor to accept a context
    dinopayResp, err := ev.dinopayClient.CreatePayment(context.Background(), &dinopayApi.Payment{
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

    err = ev.esDB.AppendEvents(dinopayEvents.BuildStreamName(payment.ID.Value.String()), outboundPaymentCreated)
    if err != nil {
        return fmt.Errorf("failed adding OutboundPaymentCreated event to the repository: %w", err)
    }

    ev.logger.Info("WithdrawalCreated event processed successfully", slog.String("withdrawal_id", withdrawalCreated.Id))

    return nil
}

func handleError(logger *slog.Logger, errMsg string, err error) error {
    logger.Error(errMsg, logattr.Error(err.Error()))
    return fmt.Errorf("%s: %w", errMsg, err)
}
