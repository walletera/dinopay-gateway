package gateway

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/events"
)

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

var _ events.Event[EventsVisitor] = OutboundPaymentCreated{}

type OutboundPaymentCreated struct {
    Id                   uuid.UUID `json:"id,omitempty"`
    WithdrawalId         uuid.UUID `json:"withdrawal_id,omitempty"`
    DinopayPaymentId     uuid.UUID `json:"dinopay_payment_id,omitempty"`
    DinopayPaymentStatus string    `json:"dinopay_payment_status,omitempty"`
    CreatedAt            int64     `json:"created_at,omitempty"`
}

func (o OutboundPaymentCreated) ID() string {
    return fmt.Sprintf("%s-%s", o.Type(), o.Id)
}

func (o OutboundPaymentCreated) Type() string {
    return "OutboundPaymentCreated"
}

func (o OutboundPaymentCreated) DataContentType() string {
    return "application/json"
}

func (o OutboundPaymentCreated) CorrelationID() string {
    panic("not implemented yet")
}

func (o OutboundPaymentCreated) Serialize() ([]byte, error) {
    data, err := json.Marshal(o)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentCreated event: %w", err)
    }
    envelope := EventEnvelope{
        Type: "OutboundPaymentCreated",
        Data: data,
    }
    return json.Marshal(envelope)
}

func (o OutboundPaymentCreated) Accept(ctx context.Context, visitor EventsVisitor) errors.ProcessingError {
    return visitor.VisitOutboundPaymentCreated(ctx, o)
}

type OutboundPaymentUpdated struct {
    Id                   uuid.UUID `json:"id,omitempty"`
    DinopayPaymentId     uuid.UUID `json:"dinopay_payment_id,omitempty"`
    DinopayPaymentStatus string    `json:"dinopay_payment_status,omitempty"`
    CreatedAt            int64     `json:"created_at,omitempty"`
}

var _ events.Event[EventsVisitor] = OutboundPaymentUpdated{}

func (o OutboundPaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", o.Type(), o.Id)
}

func (o OutboundPaymentUpdated) Type() string {
    return "OutboundPaymentUpdated"
}

func (o OutboundPaymentUpdated) DataContentType() string {
    return "application/json"
}

func (o OutboundPaymentUpdated) CorrelationID() string {
    panic("not implemented yet")
}

func (o OutboundPaymentUpdated) Serialize() ([]byte, error) {
    data, err := json.Marshal(o)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentUpdated event: %w", err)
    }
    envelope := EventEnvelope{
        Type: "OutboundPaymentUpdated",
        Data: data,
    }
    return json.Marshal(envelope)
}

func (o OutboundPaymentUpdated) Accept(ctx context.Context, visitor EventsVisitor) errors.ProcessingError {
    return visitor.VisitOutboundPaymentUpdated(ctx, o)
}
