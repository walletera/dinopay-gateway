package payments

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/events"
)

var _ events.Event[EventsHandler] = WithdrawalCreatedEvent{}

type WithdrawalCreatedEvent struct {
    Id          string                `json:"id"`
    UserId      string                `json:"user_id"`
    PspId       string                `json:"psp_id"`
    ExternalId  string                `json:"external_id"`
    Amount      float64               `json:"amount"`
    Currency    string                `json:"currency"`
    Status      string                `json:"status"`
    Beneficiary WithdrawalBeneficiary `json:"beneficiary"`

    correlationID string
}

type WithdrawalBeneficiary struct {
    Id          string             `json:"id"`
    Description string             `json:"description"`
    Account     BeneficiaryAccount `json:"account"`
}

type BeneficiaryAccount struct {
    Holder     string `json:"holder"`
    Number     int    `json:"number"`
    RoutingKey string `json:"routing_key"`
}

func (w WithdrawalCreatedEvent) Accept(ctx context.Context, handler EventsHandler) errors.ProcessingError {
    return handler.HandleWithdrawalCreated(ctx, w)
}

func (w WithdrawalCreatedEvent) ID() string {
    return fmt.Sprintf("%s-%s", w.Type(), w.Id)
}

func (w WithdrawalCreatedEvent) Type() string {
    return "WithdrawalCreated"
}

func (w WithdrawalCreatedEvent) CorrelationID() string {
    return w.correlationID
}

func (w WithdrawalCreatedEvent) DataContentType() string {
    return "application/json"
}

func (w WithdrawalCreatedEvent) Serialize() ([]byte, error) {
    return json.Marshal(w)
}
