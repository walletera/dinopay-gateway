package eventstoredb

import (
    "context"
    "errors"
    "io"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/eventskit/eventsourcing"
    "github.com/walletera/werrors"
)

// TODO check if there is a way to read all events from a stream
// maybe reading page by page
const readEventsMaxCount = 1000000

type DB struct {
    client *esdb.Client
}

func NewDB(client *esdb.Client) *DB {
    return &DB{
        client: client,
    }
}

func (db *DB) AppendEvents(ctx context.Context, streamName string, expectedAggregateVersion eventsourcing.ExpectedAggregateVersion, events ...events.EventData) werrors.WError {
    var eventsData []esdb.EventData
    for _, event := range events {
        data, err := event.Serialize()
        if err != nil {
            return werrors.NewNonRetryableInternalError(err.Error())
        }
        eventData := esdb.EventData{
            ContentType: esdb.ContentTypeJson,
            EventType:   event.Type(),
            Data:        data,
        }
        eventsData = append(eventsData, eventData)
    }
    var expectedRevision esdb.ExpectedRevision
    if expectedAggregateVersion.IsNew {
        expectedRevision = esdb.NoStream{}
    } else {
        expectedRevision = esdb.Revision(expectedAggregateVersion.Version)
    }
    opts := esdb.AppendToStreamOptions{
        ExpectedRevision: expectedRevision,
        Authenticated:    nil,
        Deadline:         nil,
        RequiresLeader:   false,
    }
    _, err := db.client.AppendToStream(ctx, streamName, opts, eventsData...)
    if err != nil {
        return mapAppendErrorToWalleteraError(err, expectedAggregateVersion)
    }
    return nil
}

func (db *DB) ReadEvents(ctx context.Context, streamName string) ([]eventsourcing.RetrievedEvent, werrors.WError) {
    stream, err := db.client.ReadStream(ctx, streamName, esdb.ReadStreamOptions{
        Direction:      esdb.Forwards,
        From:           esdb.Start{},
        ResolveLinkTos: true,
    }, readEventsMaxCount)

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
    esdbError, _ := esdb.FromError(err)
    switch esdbError.Code() {
    case esdb.ErrorCodeWrongExpectedVersion:
        if version.IsNew {
            return werrors.NewResourceAlreadyExistError(err.Error())
        } else {
            return werrors.NewWrongResourceVersionError(err.Error())
        }
    case esdb.ErrorCodeResourceNotFound:
        return werrors.NewResourceNotFoundError(err.Error())
    default:
        return werrors.NewRetryableInternalError(err.Error())
    }
}

func mapReadErrorToWalleteraError(err error) werrors.WError {
    esdbError, _ := esdb.FromError(err)
    switch esdbError.Code() {
    case esdb.ErrorCodeResourceNotFound:
        return werrors.NewResourceNotFoundError(err.Error())
    default:
        return werrors.NewRetryableInternalError(err.Error())
    }
}
