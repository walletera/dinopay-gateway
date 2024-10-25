package payments

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/walletera/message-processor/events"
)

type EventsDeserializer struct {
}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (d *EventsDeserializer) Deserialize(rawPayload []byte) (events.Event[EventsHandler], error) {
    var event events.EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "WithdrawalCreated":
        var withdrawalCreated WithdrawalCreatedEvent
        err := json.Unmarshal(event.Data, &withdrawalCreated)
        if err != nil {
            log.Printf("error deserializing WithdrawalCreatedEvent event data %s: %s", event.Data, err.Error())
        }
        withdrawalCreated.correlationID = event.CorrelationID
        return withdrawalCreated, nil
    default:
        log.Printf("unexpected event type: %s", event.Type)
        return nil, nil
    }
}
