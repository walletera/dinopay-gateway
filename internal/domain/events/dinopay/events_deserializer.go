package dinopay

import (
    "encoding/json"
    "fmt"

    "github.com/walletera/eventskit/events"
)

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

type EventsDeserializer struct{}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (e EventsDeserializer) Deserialize(rawEvent []byte) (events.Event[EventsHandler], error) {
    var eventEnvelope EventEnvelope
    err := json.Unmarshal(rawEvent, &eventEnvelope)
    if err != nil {
        return nil, fmt.Errorf("failed unmarshalling event envelope: %w", err)
    }
    switch eventEnvelope.Type {
    case "PaymentCreated":
        var paymentData PaymentData
        err := json.Unmarshal(eventEnvelope.Data, &paymentData)
        if err != nil {
            return nil, fmt.Errorf("failed unmarshalling PaymentCreated event: %w", err)
        }
        paymentCreated := PaymentCreated{
            EventType: "PaymentCreated",
            Data: PaymentData{
                Id:       paymentData.Id,
                Amount:   paymentData.Amount,
                Currency: paymentData.Currency,
                SourceAccount: Account{
                    AccountHolder: paymentData.SourceAccount.AccountHolder,
                    AccountNumber: paymentData.SourceAccount.AccountNumber,
                },
                DestinationAccount: Account{
                    AccountHolder: paymentData.DestinationAccount.AccountHolder,
                    AccountNumber: paymentData.DestinationAccount.AccountNumber,
                },
            },
        }
        return paymentCreated, nil
    default:
        return nil, fmt.Errorf("unexpected event type: %s", eventEnvelope.Type)
    }
}
