package events

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/payments-types/pkg/wogen"
    api "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

var _ events.Event[Handler] = PaymentCreated{}

const (
    PaymentCreatedType = "PaymentCreated"
)

type PaymentCreated struct {
    Id                    uuid.UUID `json:"id"`
    EventType             string    `json:"type"`
    EventAggregateVersion uint64    `json:"aggregateVersion"`
    EventCorrelationId    string    `json:"eventCorrelationId"`
    EventCreatedAt        time.Time `json:"createdAt"`

    SerializableData json.Marshaler `json:"data"`
    Data             api.Payment    `json:"-"`
}

func NewPaymentCreated(correlationId string, data api.Payment) PaymentCreated {
    return PaymentCreated{
        Id:                    uuid.New(),
        EventType:             PaymentCreatedType,
        EventAggregateVersion: 0,
        EventCorrelationId:    correlationId,
        EventCreatedAt:        time.Now(),
        SerializableData:      wogen.NewSerializationWrapper(&data),
        Data:                  data,
    }
}

func PaymentCreatedFromEnvelope(eventEnvelope events.EventEnvelope, data api.Payment) PaymentCreated {
    return PaymentCreated{
        Id:                    eventEnvelope.Id,
        EventType:             eventEnvelope.Type,
        EventAggregateVersion: eventEnvelope.AggregateVersion,
        EventCorrelationId:    eventEnvelope.CorrelationId,
        EventCreatedAt:        eventEnvelope.CreatedAt,
        SerializableData:      wogen.NewSerializationWrapper(&data),
        Data:                  data,
    }
}

func (p PaymentCreated) Accept(ctx context.Context, handler Handler) werrors.WError {
    return handler.HandlePaymentCreated(ctx, p)
}

func (p PaymentCreated) ID() string {
    return fmt.Sprintf("%s-%s", p.Type(), p.Id)
}

func (p PaymentCreated) AggregateVersion() uint64 { return p.EventAggregateVersion }

func (p PaymentCreated) Type() string {
    return PaymentCreatedType
}

func (p PaymentCreated) CorrelationID() string {
    return p.EventCorrelationId
}

func (p PaymentCreated) DataContentType() string {
    return "application/json"
}

func (p PaymentCreated) CreatedAt() time.Time { return p.EventCreatedAt }

func (p PaymentCreated) Serialize() ([]byte, error) {
    return json.Marshal(p)
}
