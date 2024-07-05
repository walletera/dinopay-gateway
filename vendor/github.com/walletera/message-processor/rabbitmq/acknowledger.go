package rabbitmq

import (
    "github.com/rabbitmq/amqp091-go"
    "github.com/walletera/message-processor/messages"
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
    return a.delivery.Nack(false, opts.Requeue)
}
