package events

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type EventEnvelope struct {
    Id               uuid.UUID `json:"id"`
    Type             string    `json:"type"`
    AggregateVersion uint64    `json:"aggregateVersion"`
    CorrelationId    string    `json:"correlationId"`
    CreatedAt        time.Time `json:"createdAt"`

    Data json.RawMessage `json:"data"`
}
