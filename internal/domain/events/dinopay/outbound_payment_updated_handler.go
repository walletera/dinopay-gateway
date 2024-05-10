package dinopay

import (
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    paymentsApi "github.com/walletera/payments/api"
)

type HandlerError struct {
    withdrawalId string
    err          error
}

func (h *HandlerError) Error() string {
    return fmt.Sprintf("error handling DinopayOutboundPaymentUpdated event for withdrawal %s", h.withdrawalId)
}

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

func (h *OutboundPaymentUpdatedHandler) Handle(outboundPaymentUpdated OutboundPaymentUpdated) *HandlerError {
    streamName := BuildStreamName(outboundPaymentUpdated.DinopayPaymentId.String())
    rawEvents, err := h.db.ReadEvents(streamName)
    if err != nil {
        return &HandlerError{
            err: fmt.Errorf("failed retrieving events from stream %s: %w", streamName, err),
        }
    }
    for _, rawEvent := range rawEvents {
        event, err := h.deserializer.Deserialize(rawEvent)
        if err != nil {
            return &HandlerError{
                err: fmt.Errorf("failed deserializing outboundPaymentUpdated event from raw event %s: %w", rawEvent, err),
            }
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
        return &HandlerError{
            err: fmt.Errorf("missing OutboundPaymentCreated event"),
        }
    }
    err := updateWithdrawalStatus(
        h.paymentsClient,
        h.outboundPaymentCreated.WithdrawalId,
        outboundPaymentUpdated.DinopayPaymentId,
        outboundPaymentUpdated.DinopayPaymentStatus,
    )
    if err != nil {
        return &HandlerError{
            withdrawalId: h.outboundPaymentCreated.WithdrawalId.String(),
            err:          err,
        }
    }
    return nil
}
