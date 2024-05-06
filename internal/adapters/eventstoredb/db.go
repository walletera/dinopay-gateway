package eventstoredb

import (
    "context"
    "errors"
    "io"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    domainEvents "github.com/walletera/dinopay-gateway/internal/domain/events"
    "github.com/walletera/message-processor/events"
)

type DB struct {
    client *esdb.Client
}

func NewDB(client *esdb.Client) *DB {
    return &DB{
        client: client,
    }
}

func (db *DB) AppendEvents(streamName string, events ...events.EventData) error {
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
    _, err := db.client.AppendToStream(context.Background(), streamName, esdb.AppendToStreamOptions{}, eventsData...)
    if err != nil {
        return err
    }
    return nil
}

func (db *DB) ReadEvents(streamName string) ([]domainEvents.RawEvent, error) {
    stream, err := db.client.ReadStream(context.Background(), streamName, esdb.ReadStreamOptions{
        Direction: esdb.Backwards,
        From:      esdb.End{},
    }, 10)

    if err != nil {
        return nil, err
    }

    defer stream.Close()

    var rawEvents []domainEvents.RawEvent
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
