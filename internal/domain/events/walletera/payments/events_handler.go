package payments

import (
    "context"
    "log/slog"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/outbound"
    "github.com/walletera/dinopay-gateway/internal/domain/ports/output/dinopay"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    dinopayapi "github.com/walletera/dinopay/api"
    "github.com/walletera/eventskit/eventsourcing"
    paymentEvents "github.com/walletera/payments-types/events"
    "github.com/walletera/werrors"
)

type EventsHandler struct {
    dinopayClient dinopay.Client
    esDB          eventsourcing.DB
    logger        *slog.Logger
}

var _ paymentEvents.Handler = (*EventsHandler)(nil)

func NewEventsHandler(dinopayClient dinopay.Client, esDB eventsourcing.DB, logger *slog.Logger, ) *EventsHandler {
    return &EventsHandler{
        dinopayClient: dinopayClient,
        esDB:          esDB,
        logger:        logger.With(logattr.Component("payments.EventsHandler")),
    }
}

func (ev *EventsHandler) HandlePaymentCreated(ctx context.Context, paymentCreated paymentEvents.PaymentCreated) werrors.WError {
    walleteraPaymentId := paymentCreated.Data.ID
    logger := ev.logger.With(
        logattr.CorrelationId(paymentCreated.CorrelationID()),
        logattr.EventType(paymentCreated.Type()),
        logattr.PaymentId(walleteraPaymentId.String()),
    )
    logger.Debug("handling PaymentCreated event")
    if !paymentCreated.Data.Debtor.AccountDetails.OneOf.IsDinopayAccountDetails() {
        return werrors.NewValidationError("invalid debtor account details type: %s", paymentCreated.Data.Beneficiary.AccountDetails.OneOf.Type)
    }
    if !paymentCreated.Data.Beneficiary.AccountDetails.OneOf.IsDinopayAccountDetails() {
        return werrors.NewValidationError("invalid beneficiary account details type: %s", paymentCreated.Data.Beneficiary.AccountDetails.OneOf.Type)
    }
    dinopayResp, err := ev.dinopayClient.CreatePayment(ctx, &dinopayapi.Payment{
        Amount:   paymentCreated.Data.Amount,
        Currency: string(paymentCreated.Data.Currency),
        SourceAccount: dinopayapi.Account{
            AccountHolder: paymentCreated.Data.Debtor.AccountDetails.OneOf.DinopayAccountDetails.AccountHolder,
            AccountNumber: paymentCreated.Data.Debtor.AccountDetails.OneOf.DinopayAccountDetails.AccountHolder,
        },
        DestinationAccount: dinopayapi.Account{
            AccountHolder: paymentCreated.Data.Beneficiary.AccountDetails.OneOf.DinopayAccountDetails.AccountHolder,
            AccountNumber: paymentCreated.Data.Beneficiary.AccountDetails.OneOf.DinopayAccountDetails.AccountNumber,
        },
        CustomerTransactionId: dinopayapi.OptString{
            Value: walleteraPaymentId.String(),
            Set:   true,
        },
    })
    if err != nil {
        werr := werrors.NewRetryableInternalError("failed creating payment on dinopay: %s", err.Error())
        logger.Error(werr.Error())
        return werr
    }
    if dinopayResp == nil {
        werr := werrors.NewRetryableInternalError("dinopay response is nil")
        logger.Error(werr.Error())
        return werr
    }
    dinopayPayment, ok := dinopayResp.(*dinopayapi.Payment)
    if !ok {
        werr := werrors.NewNonRetryableInternalError("unexpected dinopay response type %t:", dinopayResp)
        logger.Error(werr.Error())
        return werr
    }

    logger.Info("dinopay dinopayPayment created successfully")

    outboundPaymentCreated := outbound.PaymentCreated{
        Id:                   uuid.New(),
        PaymentId:            walleteraPaymentId,
        DinopayPaymentId:     dinopayPayment.ID.Value,
        DinopayPaymentStatus: string(dinopayPayment.Status.Value),
    }

    streamName := outbound.BuildOutboundPaymentStreamName(dinopayPayment.ID.Value.String())

    appendEventsErr := ev.esDB.AppendEvents(
        ctx,
        streamName,
        eventsourcing.ExpectedAggregateVersion{IsNew: true},
        outboundPaymentCreated,
    )
    if appendEventsErr != nil {
        werr := werrors.NewWrappedError(
            appendEventsErr,
            "failed handling outbound PaymentCreated event",
            streamName,
        )
        logger.Error(werr.Error())
        return werr
    }

    logger.Info("PaymentCreated event processed successfully")

    return nil
}

func (ev *EventsHandler) HandlePaymentUpdated(_ context.Context, _ paymentEvents.PaymentUpdated) werrors.WError {
    // Ignore, nothing to do
    return nil
}
