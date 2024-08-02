package events

import "github.com/walletera/message-processor/events"

type RawEvent []byte

type DB interface {
    // TODO add ctx as first parameter
    AppendEvents(streamName string, event ...events.EventData) error
    ReadEvents(streamName string) ([]RawEvent, error)
}
