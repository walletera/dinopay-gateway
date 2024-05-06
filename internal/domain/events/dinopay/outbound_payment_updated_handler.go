package dinopay

import (
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    paymentsApi "github.com/walletera/payments/api"
)

type OutboundPaymentUpdatedHandler struct {
    db             events.DB
    paymentsClient *paymentsApi.Client
    deserializer   *EventsDeserializer

    outboundPaymentCreated *OutboundPaymentCreated
}

func NewOutboundPaymentUpdatedHandler(db events.DB, client *paymentsApi.Client) *OutboundPaymentUpdatedHandler {
    return &OutboundPaymentUpdatedHandler{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
    }
}

func (h *OutboundPaymentUpdatedHandler) Handle(outboundPaymentUpdated OutboundPaymentUpdated) error {
    streamName := BuildStreamName(outboundPaymentUpdated.DinopayPaymentId.String())
    rawEvents, err := h.db.ReadEvents(streamName)
    if err != nil {
        return fmt.Errorf("failed retrieving events from stream %s: %w", streamName, err)
    }
    for _, rawEvent := range rawEvents {
        event, err := h.deserializer.Deserialize(rawEvent)
        if err != nil {
            return fmt.Errorf("failed deserializing outboundPaymentUpdated event from raw event %s: %w", rawEvent, err)
        }
        err = event.Accept(h)
        if err != nil {

        }
    }
    return nil
}

func (h *OutboundPaymentUpdatedHandler) VisitOutboundPaymentCreated(outboundPaymentCreated OutboundPaymentCreated) error {
    h.outboundPaymentCreated = &outboundPaymentCreated
    return nil
}

func (h *OutboundPaymentUpdatedHandler) VisitOutboundPaymentUpdated(outboundPaymentUpdated OutboundPaymentUpdated) error {
    if h.outboundPaymentCreated == nil {
        return fmt.Errorf("can't handle OutboundPaymentUpdated for DinoPay payment %s"+
            " event because the OutboundPaymentCreated event is missing", outboundPaymentUpdated.DinopayPaymentId)
    }
    return updateWithdrawalStatus(
        h.paymentsClient,
        h.outboundPaymentCreated.WithdrawalId,
        outboundPaymentUpdated.DinopayPaymentId,
        outboundPaymentUpdated.DinopayPaymentStatus,
    )
    return nil
}
