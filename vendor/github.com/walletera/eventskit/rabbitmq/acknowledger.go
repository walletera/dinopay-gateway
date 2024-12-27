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
    // TODO implement a way to limit the retries
    return a.delivery.Nack(false, false)
}
