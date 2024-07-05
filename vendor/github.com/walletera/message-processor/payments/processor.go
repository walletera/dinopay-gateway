package payments

import "github.com/walletera/message-processor/messages"

// NewProcessor returns a messages.Processor which is specific to payments events
func NewProcessor(messageConsumer messages.Consumer, eventsVisitor EventsVisitor, opts ...messages.ProcessorOpt) *messages.Processor[EventsVisitor] {
    return messages.NewProcessor[EventsVisitor](
        messageConsumer,
        NewEventsDeserializer(),
        eventsVisitor,
        opts...,
    )
}
