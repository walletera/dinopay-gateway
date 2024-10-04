package inbound

import (
    "context"
    "log/slog"

    "github.com/walletera/dinopay-gateway/pkg/logattr"
    procerrors "github.com/walletera/message-processor/errors"
    paymentsapi "github.com/walletera/payments-types/api"
)

type EventsHandler interface {
    HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) procerrors.ProcessingError
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

func (ev *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) procerrors.ProcessingError {

    depositPostBody := &paymentsapi.Payment{
        ID: paymentsapi.NewOptUUID(inboundPaymentReceived.DepositId),
        // FIXME
        Amount:     float64(inboundPaymentReceived.Amount),
        Currency:   inboundPaymentReceived.Currency,
        CustomerId: paymentsapi.NewOptUUID(inboundPaymentReceived.CustomerId),
        ExternalId: paymentsapi.NewOptUUID(inboundPaymentReceived.DinopayPaymentId),
    }
    _, err := ev.paymentsApiClient.PostPayment(ctx, depositPostBody, paymentsapi.PostPaymentParams{})
    if err != nil {
        // TODO this may be a retryable error
        ev.logger.Error("failed creating deposit on payments api", logattr.Error(err.Error()))
        return procerrors.NewInternalError(err.Error())
    }
    // TODO handle response
    ev.logger.Info("Gateway event InboundPaymentReceived processed successfully", logattr.EventType(inboundPaymentReceived.Type()))
    return nil
}
