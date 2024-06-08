package gateway

import (
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
)

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

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

func (o OutboundPaymentCreated) Accept(visitor EventsVisitor) error {
    return visitor.VisitOutboundPaymentCreated(o)
}

type OutboundPaymentUpdated struct {
    Id                   uuid.UUID `json:"id,omitempty"`
    DinopayPaymentId     uuid.UUID `json:"dinopay_payment_id,omitempty"`
    DinopayPaymentStatus string    `json:"dinopay_payment_status,omitempty"`
    CreatedAt            int64     `json:"created_at,omitempty"`
}

func (o OutboundPaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", o.Type(), o.Id)
}

func (o OutboundPaymentUpdated) Type() string {
    return "OutboundPaymentUpdated"
}

func (o OutboundPaymentUpdated) DataContentType() string {
    return "application/json"
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

func (o OutboundPaymentUpdated) Accept(visitor EventsVisitor) error {
    return visitor.VisitOutboundPaymentUpdated(o)
}
