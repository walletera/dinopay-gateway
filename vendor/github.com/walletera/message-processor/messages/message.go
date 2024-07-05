package messages

import (
    "github.com/walletera/message-processor/errors"
)

type NackOpts struct {
    Requeue      bool
    ErrorCode    errors.ErrorCode
    ErrorMessage string
}

type Acknowledger interface {
    Ack() error
    Nack(opts NackOpts) error
}

type Message struct {
    payload      []byte
    acknowledger Acknowledger
}

func NewMessage(payload []byte, acknowledger Acknowledger) Message {
    return Message{payload: payload, acknowledger: acknowledger}
}

func (m Message) Payload() []byte {
    return m.payload
}

func (m Message) Acknowledger() Acknowledger {
    return m.acknowledger
}
