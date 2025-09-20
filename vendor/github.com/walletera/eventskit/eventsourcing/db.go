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

// DB defines the interface for appending and reading events in an event store.
type DB interface {
    // AppendEvents appends one or more events to an event stream with optimistic concurrency control.
    //
    // Parameters:
    //   - ctx: The context for the operation, which can be used for cancellation or timeout.
    //   - streamName: The name of the stream to append events to.
    //   - resourceVersion: The expected version of the stream for optimistic concurrency control.
    //     If IsNew is true, the stream is expected not to exist yet.
    //     If IsNew is false, the Version field specifies the expected current version of the stream.
    //   - event: One or more events to append to the stream.
    //
    // Returns:
    //   - uint64: The next expected version to use for subsequent appends to this stream.
    //   - werrors.WError: An error if the operation fails, which can be:
    //     * ResourceAlreadyExistError: If IsNew is true but the stream already exists.
    //     * WrongResourceVersionError: If the specified Version doesn't match the current stream version.
    //     * ResourceNotFoundError: If the stream doesn't exist but was expected to.
    //     * Other errors for system or connection issues.
    AppendEvents(ctx context.Context, streamName string, resourceVersion ExpectedAggregateVersion, event ...events.EventData) (uint64, werrors.WError)

    // ReadEvents retrieves all events from a specified event stream.
    //
    // Parameters:
    //   - ctx: The context for the operation, which can be used for cancellation or timeout.
    //   - streamName: The name of the stream to read events from.
    //
    // Returns:
    //   - []RetrievedEvent: A slice of retrieved events, each containing the raw event data and its aggregate version.
    //   - werrors.WError: An error if the operation fails, which can be:
    //     * ResourceNotFoundError: If the specified stream doesn't exist.
    //     * Other errors for system or connection issues.
    ReadEvents(ctx context.Context, streamName string) ([]RetrievedEvent, werrors.WError)
}
