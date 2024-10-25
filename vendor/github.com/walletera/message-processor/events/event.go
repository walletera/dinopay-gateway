package events

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type Event[Handler any] interface {
    EventData

    Accept(ctx context.Context, Handler Handler) errors.ProcessingError
}
