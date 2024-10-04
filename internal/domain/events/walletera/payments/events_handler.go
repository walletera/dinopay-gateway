package payments

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/outbound"
    "github.com/walletera/dinopay-gateway/internal/domain/ports/output/dinopay"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    dinopayapi "github.com/walletera/dinopay/api"
    "github.com/walletera/message-processor/errors"
    paymentEvents "github.com/walletera/payments-types/events"
)

type EventsHandler struct {
    dinopayClient dinopay.Client
    esDB          events.DB
    logger        *slog.Logger
}

var _ paymentEvents.Handler = (*EventsHandler)(nil)

func NewEventsHandler(dinopayClient dinopay.Client, esDB events.DB, logger *slog.Logger, ) *EventsHandler {
    return &EventsHandler{
        dinopayClient: dinopayClient,
        esDB:          esDB,
        logger:        logger.With(logattr.Component("payments.EventsHandler")),
    }
}

func (ev *EventsHandler) HandlePaymentCreated(ctx context.Context, paymentCreated paymentEvents.PaymentCreated) errors.ProcessingError {
    paymentCreatedId := paymentCreated.Data.ID.Value.String()
    logger := ev.logger.With(
        logattr.EventType("paymentCreated"),
        logattr.WithdrawalId(paymentCreatedId),
    )

    withdrawalId, err := uuid.Parse(paymentCreatedId)
    if err != nil {
        handleError(logger, "failed parsing paymentCreated uuid", err)
    }

    dinopayResp, err := ev.dinopayClient.CreatePayment(ctx, &dinopayapi.Payment{
        Amount:   paymentCreated.Data.Amount,
        Currency: paymentCreated.Data.Currency,
        SourceAccount: dinopayapi.Account{
            AccountHolder: "hardcodedSourceAccountHolder",
            AccountNumber: "hardcodedSourceAccountNumber",
        },
        DestinationAccount: dinopayapi.Account{
            AccountHolder: paymentCreated.Data.Beneficiary.Value.AccountHolder.Value,
            AccountNumber: paymentCreated.Data.Beneficiary.Value.AccountNumber.Value,
        },
        CustomerTransactionId: dinopayapi.OptString{
            Value: paymentCreatedId,
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

    outboundPaymentCreated := outbound.PaymentCreated{
        Id:                   uuid.New(),
        WithdrawalId:         withdrawalId,
        DinopayPaymentId:     payment.ID.Value,
        DinopayPaymentStatus: string(payment.Status.Value),
    }

    err = ev.esDB.AppendEvents(outbound.BuildOutboundPaymentStreamName(payment.ID.Value.String()), outboundPaymentCreated)
    if err != nil {
        return errors.NewInternalError(fmt.Sprintf("failed adding OutboundPaymentCreated event to the repository: %s", err.Error()))
    }

    logger.Info("PaymentCreated event processed successfully")

    return nil
}

func (ev *EventsHandler) HandlePaymentUpdated(ctx context.Context, paymentCreatedEvent paymentEvents.PaymentUpdated) errors.ProcessingError {
    //TODO implement me
    panic("implement me")
}

func handleError(logger *slog.Logger, errMsg string, err error) errors.ProcessingError {
    logger.Error(errMsg, logattr.Error(err.Error()))
    return errors.NewInternalError(fmt.Sprintf("%s: %s", errMsg, err.Error()))
}
