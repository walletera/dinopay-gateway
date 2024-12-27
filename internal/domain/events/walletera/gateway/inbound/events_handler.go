package inbound

import (
    "context"
    "log/slog"

    "github.com/walletera/dinopay-gateway/pkg/logattr"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

type EventsHandler interface {
    HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError
}

type EventsHandlerImpl struct {
    paymentsApiClient *paymentsapi.Client
    deserializer      *EventsDeserializer
    logger            *slog.Logger
}

func NewEventsHandlerImpl(client *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        paymentsApiClient: client,
        deserializer:      NewEventsDeserializer(),
        logger:            logger.With(logattr.Component("gateway.inbound.EventsHandlerImpl")),
    }
}

func (ev *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError {

    depositPostBody := &paymentsapi.Payment{
        ID: inboundPaymentReceived.PaymentId,
        // FIXME
        Amount:     float64(inboundPaymentReceived.Amount),
        Currency:   inboundPaymentReceived.Currency,
        CustomerId: paymentsapi.NewOptUUID(inboundPaymentReceived.CustomerId),
        ExternalId: paymentsapi.NewOptUUID(inboundPaymentReceived.DinopayPaymentId),
    }
    _, err := ev.paymentsApiClient.PostPayment(ctx, depositPostBody, paymentsapi.PostPaymentParams{})
    if err != nil {
        // TODO handle this error properly
        ev.logger.Error("failed creating deposit on payments api", logattr.Error(err.Error()))
        return werrors.NewRetryableInternalError(err.Error())
    }
    // TODO handle response
    ev.logger.Info("Gateway event InboundPaymentReceived processed successfully", logattr.EventType(inboundPaymentReceived.Type()))
    return nil
}
