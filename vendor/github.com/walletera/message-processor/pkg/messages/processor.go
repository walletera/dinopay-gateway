package messages

import (
    "fmt"
    "github.com/walletera/message-processor/pkg/events"
    "log"
)

type Processor[Visitor any] struct {
    messageConsumer    Consumer
    eventsDeserializer events.Deserializer[Visitor]
    eventsVisitor      Visitor
}

func NewProcessor[Visitor any](
    messageConsumer Consumer,
    eventsDeserializer events.Deserializer[Visitor],
    eventsVisitor Visitor,
) *Processor[Visitor] {

    return &Processor[Visitor]{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsVisitor:      eventsVisitor,
    }
}

func (p *Processor[Visitor]) Start() error {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go p.processMsgs(msgCh)
    return nil
}

func (p *Processor[Visitor]) processMsgs(ch <-chan Message) {
    for msg := range ch {
        go p.processMsg(msg)
    }
}

func (p *Processor[Visitor]) processMsg(message Message) {
    event, err := p.eventsDeserializer.Deserialize(message.Payload)
    if err != nil {
        p.printErrorLog(message, err)
        return
    }
    if event != nil {
        err = event.Accept(p.eventsVisitor)
        if err != nil {
            p.printErrorLog(message, err)
        }
    }
}

func (p *Processor[V]) printErrorLog(msg Message, err error) {
    log.Printf("error processing message with payload: %s: %s", msg.Payload, err.Error())
}
