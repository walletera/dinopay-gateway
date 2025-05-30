package events

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/payments-types/pkg/wogen"
    "github.com/walletera/payments-types/pkg/wuuid"
    api "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

var _ events.Event[Handler] = PaymentUpdated{}

const (
    PaymentUpdatedType = "PaymentUpdated"
)

type PaymentUpdated struct {
    Id               uuid.UUID      `json:"id"`
    EventType        string         `json:"type"`
    CorrelationId    string         `json:"correlationId"`
    SerializableData json.Marshaler `json:"data"`

    Data api.PaymentUpdate `json:"-"`
}

func NewPaymentUpdated(correlationId string, data api.PaymentUpdate) PaymentUpdated {
    return PaymentUpdated{
        Id:               wuuid.NewUUID(),
        EventType:        PaymentUpdatedType,
        CorrelationId:    correlationId,
        SerializableData: wogen.NewSerializationWrapper(&data),
        Data:             data,
    }
}

func (w PaymentUpdated) Accept(ctx context.Context, visitor Handler) werrors.WError {
    return visitor.HandlePaymentUpdated(ctx, w)
}

func (w PaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", w.Type(), w.Id)
}

func (w PaymentUpdated) Type() string {
    return PaymentUpdatedType
}

func (w PaymentUpdated) CorrelationID() string {
    return w.CorrelationId
}

func (w PaymentUpdated) DataContentType() string {
    return "application/json"
}

func (w PaymentUpdated) Serialize() ([]byte, error) {
    return json.Marshal(w)
}
