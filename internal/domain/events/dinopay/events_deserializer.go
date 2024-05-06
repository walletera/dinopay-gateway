package dinopay

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

func (d *EventsDeserializer) Deserialize(rawPayload []byte) (events.Event[EventsVisitor], error) {
    var event EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "OutboundPaymentCreated":
        var outboundPaymentCreated OutboundPaymentCreated
        err := json.Unmarshal(event.Data, &outboundPaymentCreated)
        if err != nil {
            log.Printf("error deserializing OutboundPaymentCreated event data %s: %s", event.Data, err.Error())
        }
        return outboundPaymentCreated, nil
    case "OutboundPaymentUpdated":
        var outboundPaymentUpdated OutboundPaymentUpdated
        err := json.Unmarshal(event.Data, &outboundPaymentUpdated)
        if err != nil {
            log.Printf("error deserializing OutboundPaymentUpdated event data %s: %s", event.Data, err.Error())
        }
        return outboundPaymentUpdated, nil
    default:
        return nil, fmt.Errorf("unexpected event type: %s", event.Type)
    }
}
