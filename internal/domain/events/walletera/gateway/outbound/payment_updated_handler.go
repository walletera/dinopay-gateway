package outbound

import (
    "context"
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    "github.com/walletera/message-processor/errors"
    paymentsApi "github.com/walletera/payments-types/api"
)

type HandlerError struct {
    withdrawalId string
    err          error
}

func (h *HandlerError) Error() string {
    return fmt.Sprintf("error handling DinopayOutboundPaymentUpdated event for withdrawal %s", h.withdrawalId)
}

type PaymentUpdatedHandler struct {
    db             events.DB
    paymentsClient *paymentsApi.Client
    deserializer   *EventsDeserializer

    outboundPaymentCreated *PaymentCreated
}

func NewOutboundPaymentUpdatedHandler(db events.DB, client *paymentsApi.Client) *PaymentUpdatedHandler {
    return &PaymentUpdatedHandler{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
    }
}

func (h *PaymentUpdatedHandler) Handle(ctx context.Context, outboundPaymentUpdated PaymentUpdated) errors.ProcessingError {
    streamName := BuildOutboundPaymentStreamName(outboundPaymentUpdated.DinopayPaymentId.String())
    rawEvents, err := h.db.ReadEvents(streamName)
    if err != nil {
        return errors.NewInternalError(fmt.Sprintf("failed retrieving events from stream %s: %s", streamName, err.Error()))
    }
    for _, rawEvent := range rawEvents {
        event, err := h.deserializer.Deserialize(rawEvent)
        if err != nil {
            return errors.NewInternalError(fmt.Sprintf("failed deserializing outboundPaymentUpdated event from raw event %s: %s", rawEvent, err.Error()))
        }
        err = event.Accept(ctx, h)
        if err != nil {

        }
    }
    return nil
}

func (h *PaymentUpdatedHandler) VisitOutboundPaymentCreated(_ context.Context, outboundPaymentCreated PaymentCreated) errors.ProcessingError {
    h.outboundPaymentCreated = &outboundPaymentCreated
    return nil
}

func (h *PaymentUpdatedHandler) VisitOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated PaymentUpdated) errors.ProcessingError {
    if h.outboundPaymentCreated == nil {
        return errors.NewInternalError("missing OutboundPaymentCreated event")
    }
    err := updatePaymentStatus(ctx, h.paymentsClient, h.outboundPaymentCreated.WithdrawalId, outboundPaymentUpdated.DinopayPaymentId, outboundPaymentUpdated.DinopayPaymentStatus)
    if err != nil {
        return errors.NewInternalError(err.Error())
    }
    return nil
}
