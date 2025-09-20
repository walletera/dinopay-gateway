package events

import "time"

type EventData interface {
    // ID returns the id of the event
    ID() string
    // Type returns the type of the event
    Type() string
    // AggregateVersion returns the version of the aggregate this event belongs to
    AggregateVersion() uint64
    // CorrelationID returns the correlation id of the event
    // The correlation id is used for tracing linked events
    CorrelationID() string
    // DataContentType returns the content type of the event data
    DataContentType() string
    // CreatedAt returns the creation time of the event
    CreatedAt() time.Time
    // Serialize serializes the event data into a byte array for transmission or storage and returns an error if serialization fails.
    Serialize() ([]byte, error)
}
