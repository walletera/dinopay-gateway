package eventstoredb

import "github.com/walletera/message-processor/messages"

type Acknowledger struct {
}

func (a Acknowledger) Ack() error {
	//TODO implement me
	return nil
}

func (a Acknowledger) Nack(opts messages.NackOpts) error {
	//TODO implement me
	return nil
}
