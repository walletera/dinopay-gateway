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

func (pc PaymentCreated) ID() string {
    return pc.Id.String()
}

func (pc PaymentCreated) Type() string {
    return pc.EventType
}

func (pc PaymentCreated) CorrelationID() string {
    //TODO implement me
    panic("implement me")
}

func (pc PaymentCreated) DataContentType() string {
    return "application/json"
}

func (pc PaymentCreated) AggregateVersion() uint64 {
    return 0
}

func (pc PaymentCreated) CreatedAt() time.Time {
    return pc.Time
}

func (pc PaymentCreated) Serialize() ([]byte, error) {
    return json.Marshal(pc)
}

func (pc PaymentCreated) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandlePaymentCreated(ctx, pc)
}
