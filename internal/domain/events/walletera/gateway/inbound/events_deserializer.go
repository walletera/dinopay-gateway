package inbound

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/eventskit/events"
)

type EventsDeserializer struct {
}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (d *EventsDeserializer) Deserialize(rawPayload []byte) (events.Event[EventsHandler], error) {
    var event gateway.EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "InboundPaymentReceived":
        var paymentReceived PaymentReceived
        err := json.Unmarshal(event.Data, &paymentReceived)
        if err != nil {
            log.Printf("error deserializing InboundPaymentReceived event data %s: %s", event.Data, err.Error())
        }
        return paymentReceived, nil
    default:
        return nil, fmt.Errorf("unexpected event type: %s", event.Type)
    }
}
