package outbound

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

type PaymentUpdated struct {
    Id                              uuid.UUID `json:"id,omitempty"`
    DinopayPaymentId                uuid.UUID `json:"dinopay_payment_id,omitempty"`
    DinopayPaymentStatus            string    `json:"dinopay_payment_status,omitempty"`
    OutboundPaymentAggregateVersion uint64    `json:"aggregate_version,omitempty"`
    EventCreatedAt                  int64     `json:"created_at,omitempty"`
}

var _ events.Event[EventsHandler] = PaymentUpdated{}

func (pu PaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", pu.Type(), pu.Id)
}

func (pu PaymentUpdated) Type() string {
    return "OutboundPaymentUpdated"
}

func (pu PaymentUpdated) DataContentType() string {
    return "application/json"
}

func (pu PaymentUpdated) CorrelationID() string {
    panic("not implemented yet")
}

func (pu PaymentUpdated) AggregateVersion() uint64 {
    return pu.OutboundPaymentAggregateVersion
}

func (pu PaymentUpdated) CreatedAt() time.Time {
    return time.UnixMilli(pu.EventCreatedAt)
}

func (pu PaymentUpdated) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleOutboundPaymentUpdated(ctx, pu)
}

func (pu PaymentUpdated) Serialize() ([]byte, error) {
    data, err := json.Marshal(pu)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentUpdated event: %w", err)
    }
    envelope := gateway.EventEnvelope{
        Type: "OutboundPaymentUpdated",
        Data: data,
    }
    return json.Marshal(envelope)
}
