package eventstoredb

import (
    "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
    "github.com/walletera/eventskit/messages"
)

type Acknowledger struct {
    persistentSubscription *kurrentdb.PersistentSubscription
    eventAppeared          *kurrentdb.EventAppeared
}

func NewAcknowledger(persistentSubscription *kurrentdb.PersistentSubscription, event *kurrentdb.EventAppeared) *Acknowledger {
    return &Acknowledger{persistentSubscription: persistentSubscription, eventAppeared: event}
}

func (a *Acknowledger) Ack() error {
    return a.persistentSubscription.Ack(a.eventAppeared.Event)
}

func (a *Acknowledger) Nack(opts messages.NackOpts) error {
    var err error
    if opts.Requeue {
        err = a.persistentSubscription.Nack(opts.ErrorMessage, kurrentdb.NackActionRetry, a.eventAppeared.Event)
    } else {
        err = a.persistentSubscription.Nack(opts.ErrorMessage, kurrentdb.NackActionPark, a.eventAppeared.Event)
    }
    return err
}
