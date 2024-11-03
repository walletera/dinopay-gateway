package eventstoredb

import (
    "context"
    "errors"
    "io"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/walletera/message-processor/events"
    "github.com/walletera/message-processor/eventsourcing"
)

type DB struct {
    client *esdb.Client
}

func NewDB(client *esdb.Client) *DB {
    return &DB{
        client: client,
    }
}

func (db *DB) AppendEvents(ctx context.Context, streamName string, events ...events.EventData) error {
    var eventsData []esdb.EventData
    for _, event := range events {
        data, err := event.Serialize()
        if err != nil {
            return err
        }
        eventData := esdb.EventData{
            ContentType: esdb.ContentTypeJson,
            EventType:   event.Type(),
            Data:        data,
        }
        eventsData = append(eventsData, eventData)
    }
    _, err := db.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{}, eventsData...)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) ReadEvents(ctx context.Context, streamName string) ([]eventsourcing.RawEvent, error) {
    stream, err := db.client.ReadStream(ctx, streamName, esdb.ReadStreamOptions{
        Direction: esdb.Backwards,
        From:      esdb.End{},
    }, 10)

    if err != nil {
        return nil, err
    }

    defer stream.Close()

    var rawEvents []eventsourcing.RawEvent
    for {
        event, err := stream.Recv()

        if errors.Is(err, io.EOF) {
            break
        }

        if err != nil {
            return nil, err
        }

        rawEvents = append(rawEvents, event.Event.Data)
    }

    return rawEvents, nil
}
