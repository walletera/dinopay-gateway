package eventstoredb

import (
    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/walletera/eventskit/messages"
)

// TODO make this configurable
const maxRetry = 3

type Acknowledger struct {
    persistentSubscription *esdb.PersistentSubscription
    eventAppeared          *esdb.EventAppeared
}

func NewAcknowledger(persistentSubscription *esdb.PersistentSubscription, event *esdb.EventAppeared) *Acknowledger {
    return &Acknowledger{persistentSubscription: persistentSubscription, eventAppeared: event}
}

func (a *Acknowledger) Ack() error {
    return a.persistentSubscription.Ack(a.eventAppeared.Event)
}

func (a *Acknowledger) Nack(opts messages.NackOpts) error {
    var err error
    if opts.Requeue && a.eventAppeared.RetryCount <= maxRetry {
        err = a.persistentSubscription.Nack(opts.ErrorMessage, esdb.NackActionRetry, a.eventAppeared.Event)
    } else {
        err = a.persistentSubscription.Nack(opts.ErrorMessage, esdb.NackActionPark, a.eventAppeared.Event)
    }
    return err
}
