package payments

import (
    "context"
    "fmt"
    "log/slog"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events"
    gatewayevents "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/dinopay-gateway/internal/domain/ports/output/dinopay"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    dinopayapi "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/payments"
)

type EventsVisitor struct {
    dinopayClient dinopay.Client
    esDB          events.DB
    logger        *slog.Logger
}

var _ payments.EventsVisitor = (*EventsVisitor)(nil)

func NewEventsVisitor(dinopayClient dinopay.Client, esDB events.DB, logger *slog.Logger, ) *EventsVisitor {
    return &EventsVisitor{
        dinopayClient: dinopayClient,
        esDB:          esDB,
        logger:        logger.With(logattr.Component("payments.EventsVisitor")),
    }
}

func (ev *EventsVisitor) VisitWithdrawalCreated(ctx context.Context, withdrawalCreated payments.WithdrawalCreatedEvent) errors.ProcessingError {
    logger := ev.logger.With(
        logattr.EventType("WithdrawalCreated"),
        logattr.WithdrawalId(withdrawalCreated.Id),
    )

    withdrawalId, err := uuid.Parse(withdrawalCreated.Id)
    if err != nil {
        handleError(logger, "failed parsing WithdrawalCreated uuid", err)
    }

    dinopayResp, err := ev.dinopayClient.CreatePayment(ctx, &dinopayapi.Payment{
        Amount:   withdrawalCreated.Amount,
        Currency: withdrawalCreated.Currency,
        SourceAccount: dinopayapi.Account{
            AccountHolder: "hardcodedSourceAccountHolder",
            AccountNumber: "hardcodedSourceAccountNumber",
        },
        DestinationAccount: dinopayapi.Account{
            AccountHolder: withdrawalCreated.Beneficiary.Account.Holder,
            AccountNumber: strconv.Itoa(withdrawalCreated.Beneficiary.Account.Number),
        },
        CustomerTransactionId: dinopayapi.OptString{
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

    payment, ok := dinopayResp.(*dinopayapi.Payment)
    if !ok {
        return handleError(logger, "unexpected dinopay response", err)
    }

    logger.Info("dinopay payment created successfully")

    outboundPaymentCreated := gatewayevents.OutboundPaymentCreated{
        Id:                   uuid.New(),
        WithdrawalId:         withdrawalId,
        DinopayPaymentId:     payment.ID.Value,
        DinopayPaymentStatus: string(payment.Status.Value),
    }

    err = ev.esDB.AppendEvents(gatewayevents.BuildStreamName(payment.ID.Value.String()), outboundPaymentCreated)
    if err != nil {
        return errors.NewInternalError(fmt.Sprintf("failed adding OutboundPaymentCreated event to the repository: %s", err.Error()))
    }

    logger.Info("WithdrawalCreated event processed successfully")

    return nil
}

func handleError(logger *slog.Logger, errMsg string, err error) errors.ProcessingError {
    logger.Error(errMsg, logattr.Error(err.Error()))
    return errors.NewInternalError(fmt.Sprintf("%s: %s", errMsg, err.Error()))
}
