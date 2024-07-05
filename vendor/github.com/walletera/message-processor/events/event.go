package events

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type EventData interface {
    ID() string
    Type() string
    CorrelationID() string
    DataContentType() string
    Serialize() ([]byte, error)
}

type Event[Visitor any] interface {
    EventData

    Accept(ctx context.Context, visitor Visitor) errors.ProcessingError
}
