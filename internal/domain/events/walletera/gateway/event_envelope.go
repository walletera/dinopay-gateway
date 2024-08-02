package gateway

import (
    "encoding/json"
)

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}
