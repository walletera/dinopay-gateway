package eventsourcing

import (
    "context"

    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

type ExpectedAggregateVersion struct {
    IsNew   bool
    Version uint64
}

type RetrievedEvent struct {
    RawEvent         []byte
    AggregateVersion uint64
}

type DB interface {
    AppendEvents(ctx context.Context, streamName string, resourceVersion ExpectedAggregateVersion, event ...events.EventData) werrors.WError
    ReadEvents(ctx context.Context, streamName string) ([]RetrievedEvent, werrors.WError)
}
