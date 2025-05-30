package dinopay

import (
    "context"
    "encoding/json"
    "time"

    "github.com/google/uuid"
    "github.com/walletera/werrors"
)

type PaymentCreated struct {
    Id        uuid.UUID   `json:"id"`
    EventType string      `json:"type"`
    Time      time.Time   `json:"time"`
    Data      PaymentData `json:"data"`
}

type PaymentData struct {
    Id                 uuid.UUID `json:"id"`
    Amount             float64   `json:"amount"`
    Currency           string    `json:"currency"`
    SourceAccount      Account   `json:"sourceAccount"`
    DestinationAccount Account   `json:"destinationAccount"`
}

type Account struct {
    AccountHolder string `json:"accountHolder"`
    AccountNumber string `json:"accountNumber"`
}

func (p PaymentCreated) ID() string {
    return p.Id.String()
}

func (p PaymentCreated) Type() string {
    return p.EventType
}

func (p PaymentCreated) CorrelationID() string {
    //TODO implement me
    panic("implement me")
}

func (p PaymentCreated) DataContentType() string {
    return "application/json"
}

func (p PaymentCreated) Serialize() ([]byte, error) {
    return json.Marshal(p)
}

func (p PaymentCreated) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandlePaymentCreated(ctx, p)
}
