package payments

import (
    "fmt"

    "github.com/walletera/message-processor/messages"
    "github.com/walletera/message-processor/rabbitmq"
)

const (
    RabbitMQExchangeName = "payments"
    RabbitMQExchangeType = "topic"
    RabbitMQRoutingKey   = "payments.events"
)

type RabbitMQProcessorOpt struct {
    ProcessorOpt        messages.ProcessorOpt
    RabbitmqConsumerOpt rabbitmq.ConsumerOpt
}

func NewRabbitMQProcessor(
    eventsVisitor EventsVisitor,
    queueName string,
    opts ...RabbitMQProcessorOpt,
) (*messages.Processor[EventsVisitor], error) {
    defaultConsumerOpts := []rabbitmq.ConsumerOpt{
        rabbitmq.WithExchangeName(RabbitMQExchangeName),
        rabbitmq.WithExchangeType(RabbitMQExchangeType),
        rabbitmq.WithConsumerRoutingKeys(RabbitMQRoutingKey),
        rabbitmq.WithQueueName(queueName),
    }
    rabbitmqConsumerOpts := append(defaultConsumerOpts, rabbitmqConsumerOpts(opts)...)
    rabbitMQClient, err := rabbitmq.NewClient(rabbitmqConsumerOpts...)
    if err != nil {
        return nil, fmt.Errorf("creating rabbitmq client: %w", err)
    }
    return NewProcessor(rabbitMQClient, eventsVisitor, processorOpts(opts)...), nil
}

func rabbitmqConsumerOpts(opts []RabbitMQProcessorOpt) []rabbitmq.ConsumerOpt {
    var rabbitmqConsumerOpts []rabbitmq.ConsumerOpt
    for _, opt := range opts {
        if opt.RabbitmqConsumerOpt != nil {
            rabbitmqConsumerOpts = append(rabbitmqConsumerOpts, opt.RabbitmqConsumerOpt)
        }
    }
    return rabbitmqConsumerOpts
}

func processorOpts(opts []RabbitMQProcessorOpt) []messages.ProcessorOpt {
    var processorOpts []messages.ProcessorOpt
    for _, opt := range opts {
        if opt.ProcessorOpt != nil {
            processorOpts = append(processorOpts, opt.ProcessorOpt)
        }
    }
    return processorOpts
}
