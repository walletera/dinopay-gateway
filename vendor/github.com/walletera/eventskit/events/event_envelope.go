package events

import "encoding/json"

type EventEnvelope struct {
    Type          string          `json:"type"`
    CorrelationID string          `json:"correlation_id"`
    Data          json.RawMessage `json:"data"`
}
