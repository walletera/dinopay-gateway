package events

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/events"
    "github.com/walletera/payments-types/api"
    "github.com/walletera/payments-types/pkg/wuuid"
)

var _ events.Event[Handler] = PaymentUpdated{}

type PaymentUpdated struct {
    Id            uuid.UUID         `json:"id"`
    CorrelationId string            `json:"correlationId"`
    Data          api.PaymentUpdate `json:"data"`
}

func NewPaymentUpdated(correlationId string, data api.PaymentUpdate) PaymentUpdated {
    return PaymentUpdated{
        Id:            wuuid.NewUUID(),
        CorrelationId: correlationId,
        Data:          data,
    }
}

func (w PaymentUpdated) Accept(ctx context.Context, visitor Handler) errors.ProcessingError {
    return visitor.HandlePaymentUpdated(ctx, w)
}

func (w PaymentUpdated) ID() string {
    return fmt.Sprintf("%s-%s", w.Type(), w.Id)
}

func (w PaymentUpdated) Type() string {
    return "PaymentUpdated"
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
