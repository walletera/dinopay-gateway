package outbound

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

type PaymentUpdated struct {
    Id                   uuid.UUID `json:"id,omitempty"`
    DinopayPaymentId     uuid.UUID `json:"dinopay_payment_id,omitempty"`
    DinopayPaymentStatus string    `json:"dinopay_payment_status,omitempty"`
    CreatedAt            int64     `json:"created_at,omitempty"`
}

var _ events.Event[EventsHandler] = PaymentUpdated{}

func (o PaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", o.Type(), o.Id)
}

func (o PaymentUpdated) Type() string {
    return "OutboundPaymentUpdated"
}

func (o PaymentUpdated) DataContentType() string {
    return "application/json"
}

func (o PaymentUpdated) CorrelationID() string {
    panic("not implemented yet")
}

func (o PaymentUpdated) Serialize() ([]byte, error) {
    data, err := json.Marshal(o)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentUpdated event: %w", err)
    }
    envelope := gateway.EventEnvelope{
        Type: "OutboundPaymentUpdated",
        Data: data,
    }
    return json.Marshal(envelope)
}

func (o PaymentUpdated) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleOutboundPaymentUpdated(ctx, o)
}
