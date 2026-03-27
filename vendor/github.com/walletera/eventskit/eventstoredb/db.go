package eventstoredb

import (
    "context"
    "errors"
    "io"

    "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/werrors"
)

const defaultReadEventsPageSize = 100

type DB struct {
    client             *kurrentdb.Client
    readEventsPageSize int
}

type Config func(*DB)

func WithReadEventsPageSize(count int) func(db *DB) {
    return func(db *DB) {
        db.readEventsPageSize = count
    }
}

func NewDB(client *kurrentdb.Client, configs ...Config) *DB {
    db := &DB{
        client:             client,
        readEventsPageSize: defaultReadEventsPageSize,
    }
    for _, config := range configs {
        config(db)
    }
    return db
}

func (db *DB) AppendEvents(ctx context.Context, streamName string, expectedAggregateVersion eventsourcing.ExpectedAggregateVersion, events ...events.EventData) (uint64, werrors.WError) {
    var eventsData []kurrentdb.EventData
    for _, event := range events {
        data, err := event.Serialize()
        if err != nil {
            return 0, werrors.NewNonRetryableInternalError(err.Error())
        }
        eventData := kurrentdb.EventData{
            ContentType: kurrentdb.ContentTypeJson,
            EventType:   event.Type(),
            Data:        data,
        }
        eventsData = append(eventsData, eventData)
    }
    var expectedRevision kurrentdb.StreamState
    if expectedAggregateVersion.IsNew {
        expectedRevision = kurrentdb.NoStream{}
    } else {
        expectedRevision = kurrentdb.Revision(expectedAggregateVersion.Version)
    }
    opts := kurrentdb.AppendToStreamOptions{
        StreamState:    expectedRevision,
        Authenticated:  nil,
        Deadline:       nil,
        RequiresLeader: false,
    }
    writeResult, err := db.client.AppendToStream(ctx, streamName, opts, eventsData...)
    if err != nil {
        return 0, mapAppendErrorToWalleteraError(err, expectedAggregateVersion)
    }
    return writeResult.NextExpectedVersion, nil
}

func (db *DB) ReadEvents(ctx context.Context, streamName string) ([]eventsourcing.RetrievedEvent, werrors.WError) {
    retrievedEventsPage, err := db.readEventsFrom(ctx, streamName, kurrentdb.Start{})
    if err != nil {
        return nil, err
    }

    allRetrievedStreamEvents := retrievedEventsPage

    for len(retrievedEventsPage) == db.readEventsPageSize {
        lastRetrievedEvent := retrievedEventsPage[len(retrievedEventsPage)-1]
        nextReadEventPosition := kurrentdb.StreamRevision{Value: lastRetrievedEvent.AggregateVersion + 1}
        retrievedEventsPage, err = db.readEventsFrom(ctx, streamName, nextReadEventPosition)
        if err != nil {
            return nil, err
        }
        allRetrievedStreamEvents = append(allRetrievedStreamEvents, retrievedEventsPage...)
    }

    return allRetrievedStreamEvents, nil
}

func (db *DB) readEventsFrom(ctx context.Context, streamName string, from kurrentdb.StreamPosition) ([]eventsourcing.RetrievedEvent, werrors.WError) {
    stream, err := db.client.ReadStream(
        ctx,
        streamName,
        kurrentdb.ReadStreamOptions{
            Direction:      kurrentdb.Forwards,
            From:           from,
            ResolveLinkTos: true,
        },
        uint64(db.readEventsPageSize),
    )

    if err != nil {
        return nil, mapReadErrorToWalleteraError(err)
    }

    defer stream.Close()

    var rawEvents []eventsourcing.RetrievedEvent
    for {
        event, err := stream.Recv()

        if errors.Is(err, io.EOF) {
            break
        }

        if err != nil {
            return nil, mapReadErrorToWalleteraError(err)
        }

        rawEvents = append(rawEvents, eventsourcing.RetrievedEvent{
            RawEvent:         event.Event.Data,
            AggregateVersion: event.OriginalStreamRevision().Value,
        })
    }

    return rawEvents, nil
}

func mapAppendErrorToWalleteraError(err error, version eventsourcing.ExpectedAggregateVersion) werrors.WError {
    esdbError, _ := kurrentdb.FromError(err)
    switch esdbError.Code() {
    case kurrentdb.ErrorCodeWrongExpectedVersion:
        if version.IsNew {
            return werrors.NewResourceAlreadyExistError(err.Error())
        } else {
            return werrors.NewWrongResourceVersionError(err.Error())
        }
    case kurrentdb.ErrorCodeResourceNotFound:
        return werrors.NewResourceNotFoundError(err.Error())
    default:
        return werrors.NewRetryableInternalError(err.Error())
    }
}

func mapReadErrorToWalleteraError(err error) werrors.WError {
    esdbError, _ := kurrentdb.FromError(err)
    switch esdbError.Code() {
    case kurrentdb.ErrorCodeResourceNotFound:
        return werrors.NewResourceNotFoundError(err.Error())
    default:
        return werrors.NewRetryableInternalError(err.Error())
    }
}
