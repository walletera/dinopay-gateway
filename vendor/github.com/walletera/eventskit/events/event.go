package events

import (
    "context"

    "github.com/walletera/werrors"
)

type Event[Handler any] interface {
    EventData

    Accept(ctx context.Context, Handler Handler) werrors.WError
}
