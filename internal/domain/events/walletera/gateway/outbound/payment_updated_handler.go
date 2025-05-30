package outbound

import (
    "context"
    "fmt"

    "github.com/walletera/eventskit/eventsourcing"
    paymentsApi "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

type HandlerError struct {
    withdrawalId string
    err          error
}

func (h *HandlerError) Error() string {
    return fmt.Sprintf("error handling DinopayOutboundPaymentUpdated event for withdrawal %s", h.withdrawalId)
}

type PaymentUpdatedHandler struct {
    db             eventsourcing.DB
    paymentsClient *paymentsApi.Client
    deserializer   *EventsDeserializer

    outboundPaymentCreated *PaymentCreated
}

func NewOutboundPaymentUpdatedHandler(db eventsourcing.DB, client *paymentsApi.Client) *PaymentUpdatedHandler {
    return &PaymentUpdatedHandler{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
    }
}

func (h *PaymentUpdatedHandler) Handle(ctx context.Context, outboundPaymentUpdated PaymentUpdated) werrors.WError {
    streamName := BuildOutboundPaymentStreamName(outboundPaymentUpdated.DinopayPaymentId.String())
    retrievedEvents, werr := h.db.ReadEvents(ctx, streamName)
    if werr != nil {
        return werrors.NewWrappedError(werr)
    }
    for _, retrievedEvent := range retrievedEvents {
        event, err := h.deserializer.Deserialize(retrievedEvent.RawEvent)
        if err != nil {
            return werrors.NewNonRetryableInternalError("failed deserializing outboundPaymentUpdated event from raw event %s: %s", retrievedEvent, err.Error())
        }
        err = event.Accept(ctx, h)
        if err != nil {

        }
    }
    return nil
}

func (h *PaymentUpdatedHandler) HandleOutboundPaymentCreated(_ context.Context, outboundPaymentCreated PaymentCreated) werrors.WError {
    h.outboundPaymentCreated = &outboundPaymentCreated
    return nil
}

func (h *PaymentUpdatedHandler) HandleOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated PaymentUpdated) werrors.WError {
    if h.outboundPaymentCreated == nil {
        return werrors.NewNonRetryableInternalError("missing OutboundPaymentCreated event")
    }
    err := updatePaymentStatus(ctx, h.paymentsClient, h.outboundPaymentCreated.PaymentId, outboundPaymentUpdated.DinopayPaymentId, outboundPaymentUpdated.DinopayPaymentStatus)
    if err != nil {
        return werrors.NewWrappedError(err, "failed handling outbound PaymentCreated event")
    }
    return nil
}
