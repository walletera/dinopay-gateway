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

var _ events.Event[Handler] = PaymentUpdated{}

const (
    PaymentUpdatedType = "PaymentUpdated"
)

type PaymentUpdated struct {
    Id                    uuid.UUID `json:"id"`
    EventType             string    `json:"type"`
    EventAggregateVersion uint64    `json:"aggregateVersion"`
    EventCorrelationId    string    `json:"eventCorrelationId"`
    EventCreatedAt        time.Time `json:"createdAt"`

    SerializableData json.Marshaler    `json:"data"`
    Data             api.PaymentUpdate `json:"-"`
}

func NewPaymentUpdated(correlationId string, data api.PaymentUpdate) PaymentUpdated {
    return PaymentUpdated{
        Id:                    uuid.New(),
        EventType:             PaymentUpdatedType,
        EventAggregateVersion: 0,
        EventCorrelationId:    correlationId,
        EventCreatedAt:        time.Now(),
        SerializableData:      wogen.NewSerializationWrapper(&data),
        Data:                  data,
    }
}

func PaymentUpdatedFromEnvelope(eventEnvelope events.EventEnvelope, data api.PaymentUpdate) PaymentUpdated {
    return PaymentUpdated{
        Id:                    eventEnvelope.Id,
        EventType:             eventEnvelope.Type,
        EventAggregateVersion: eventEnvelope.AggregateVersion,
        EventCorrelationId:    eventEnvelope.CorrelationId,
        EventCreatedAt:        eventEnvelope.CreatedAt,
        SerializableData:      wogen.NewSerializationWrapper(&data),
        Data:                  data,
    }
}

func (w PaymentUpdated) Accept(ctx context.Context, visitor Handler) werrors.WError {
    return visitor.HandlePaymentUpdated(ctx, w)
}

func (w PaymentUpdated) AggregateVersion() uint64 { return w.EventAggregateVersion }

func (w PaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", w.Type(), w.Id)
}

func (w PaymentUpdated) Type() string {
    return PaymentUpdatedType
}

func (w PaymentUpdated) CorrelationID() string {
    return w.EventCorrelationId
}

func (w PaymentUpdated) DataContentType() string {
    return "application/json"
}

func (w PaymentUpdated) CreatedAt() time.Time { return w.EventCreatedAt }

func (w PaymentUpdated) Serialize() ([]byte, error) {
    return json.Marshal(w)
}
