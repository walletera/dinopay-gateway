package payments

import "encoding/json"

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

type BeneficiaryAccount struct {
    Holder     string `json:"holder"`
    Number     int    `json:"number"`
    RoutingKey string `json:"routing_key"`
}

type WithdrawalBeneficiary struct {
    Id          string             `json:"id"`
    Description string             `json:"description"`
    Account     BeneficiaryAccount `json:"account"`
}

type WithdrawalCreated struct {
    Id          string                `json:"id"`
    UserId      string                `json:"user_id"`
    PspId       string                `json:"psp_id"`
    ExternalId  string                `json:"external_id"`
    Amount      float64               `json:"amount"`
    Currency    string                `json:"currency"`
    Status      string                `json:"status"`
    Beneficiary WithdrawalBeneficiary `json:"beneficiary"`
}

func (w WithdrawalCreated) Accept(visitor EventsVisitor) error {
    return visitor.VisitWithdrawalCreated(w)
}
