package inbound

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

var _ events.Event[EventsHandler] = PaymentReceived{}

type PaymentReceived struct {
    Id                 uuid.UUID `json:"id,omitempty"`
    DinopayPaymentId   uuid.UUID `json:"externalId,omitempty"`
    CustomerId         uuid.UUID `json:"customerId,omitempty"`
    PaymentId          uuid.UUID `json:"depositId,omitempty"`
    Amount             float64   `json:"amount"`
    Currency           string    `json:"currency"`
    SourceAccount      Account   `json:"sourceAccount"`
    DestinationAccount Account   `json:"destinationAccount"`
    CreatedAt          time.Time `json:"createdAt,omitempty"`
}

type Account struct {
    AccountHolder string `json:"accountHolder"`
    AccountNumber string `json:"accountNumber"`
}

func (i PaymentReceived) ID() string {
    return i.Id.String()
}

func (i PaymentReceived) Type() string {
    return "InboundPaymentReceived"
}

func (i PaymentReceived) DataContentType() string {
    return "application/json"
}

func (i PaymentReceived) CorrelationID() string {
    panic("not implemented yet")
}

func (i PaymentReceived) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleInboundPaymentReceived(ctx, i)
}

func (i PaymentReceived) Serialize() ([]byte, error) {
    data, err := json.Marshal(i)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentCreated event: %w", err)
    }
    envelope := gateway.EventEnvelope{
        Type: "InboundPaymentReceived",
        Data: data,
    }
    return json.Marshal(envelope)
}
