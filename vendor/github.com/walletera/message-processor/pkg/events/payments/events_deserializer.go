package payments

import (
    "encoding/json"
    "fmt"
    "github.com/walletera/message-processor/pkg/events"
    "log"
)

type EventsDeserializer struct {
}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (d *EventsDeserializer) Deserialize(rawPayload []byte) (events.Event[EventsVisitor], error) {
    var event EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "WithdrawalCreated":
        var withdrawalCreated WithdrawalCreated
        err := json.Unmarshal(event.Data, &withdrawalCreated)
        if err != nil {
            log.Printf("error deserializing WithdrawalCreated event data %s: %s", event.Data, err.Error())
        }
        return withdrawalCreated, nil
    default:
        log.Printf("unexpected event type: %s", event.Type)
        return nil, nil
    }
}
