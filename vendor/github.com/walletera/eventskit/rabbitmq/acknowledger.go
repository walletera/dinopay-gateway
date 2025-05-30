package rabbitmq

import (
    "github.com/rabbitmq/amqp091-go"
    "github.com/walletera/eventskit/messages"
)

type Acknowledger struct {
    delivery amqp091.Delivery
}

func NewAcknowledger(delivery amqp091.Delivery) *Acknowledger {
    return &Acknowledger{
        delivery: delivery,
    }
}

func (a *Acknowledger) Ack() error {
    return a.delivery.Ack(false)
}

func (a *Acknowledger) Nack(opts messages.NackOpts) error {
    if a.delivery.Headers == nil {
        a.delivery.Headers = make(amqp091.Table)
    }
    var deliveryCount int
    header, ok := a.delivery.Headers["w-delivery-count"]
    if ok {
        deliveryCount = header.(int)
    }
    var requeue bool
    if opts.Requeue && deliveryCount < opts.MaxRetries {
        requeue = true
    }
    a.delivery.Headers["w-delivery-count"] = deliveryCount + 1
    return a.delivery.Nack(false, requeue)
}
