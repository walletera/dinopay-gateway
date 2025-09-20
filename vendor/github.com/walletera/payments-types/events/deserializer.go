package events

import (
    "encoding/json"
    "fmt"
    "log/slog"

    "github.com/walletera/eventskit/events"
    "github.com/walletera/payments-types/privateapi"
)

// interface compliance verification
var _ events.Deserializer[Handler] = (*Deserializer)(nil)

type Deserializer struct {
    logger *slog.Logger
}

func NewDeserializer(logger *slog.Logger) *Deserializer {
    return &Deserializer{logger: logger}
}

func (d *Deserializer) Deserialize(rawPayload []byte) (events.Event[Handler], error) {
    var eventEnvelope events.EventEnvelope
    err := json.Unmarshal(rawPayload, &eventEnvelope)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch eventEnvelope.Type {
    case PaymentCreatedType:
        var payment privateapi.Payment
        err := json.Unmarshal(eventEnvelope.Data, &payment)
        if err != nil {
            return nil, fmt.Errorf("error deserializing PaymentCreated eventEnvelope data %s: %s", eventEnvelope.Data, err.Error())
        }
        return PaymentCreatedFromEnvelope(eventEnvelope, payment), nil
    case PaymentUpdatedType:
        var payment privateapi.PaymentUpdate
        err := json.Unmarshal(eventEnvelope.Data, &payment)
        if err != nil {
            return nil, fmt.Errorf("error deserializing PaymentUpdated eventEnvelope data %s: %s", eventEnvelope.Data, err.Error())
        }
        return PaymentUpdatedFromEnvelope(eventEnvelope, payment), nil
    default:
        d.logger.Warn("unexpected eventEnvelope type", slog.String("eventType", eventEnvelope.Type))
        return nil, nil
    }
}
