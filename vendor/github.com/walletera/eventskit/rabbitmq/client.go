package rabbitmq

import (
    "context"
    "fmt"

    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/eventskit/messages"
)

const (
    DefaultHost     = "localhost"
    DefaultPort     = 5672
    DefaultUser     = "guest"
    DefaultPassword = "guest"

    ManagementUIPort = 15672

    ExchangeTypeDirect = "direct"
    ExchangeTypeTopic  = "topic"
    ExchangeTypeFanout = "fanout"
)

type Client struct {
    conn        *amqp.Connection
    connChannel *amqp.Channel
    queue       amqp.Queue

    host     string
    port     uint
    user     string
    password string

    useDefaultExchange  bool
    exchangeName        string
    exchangeType        string
    queueName           string
    consumerRoutingKeys []string
}

type ConsumerOpt func(consumer *Client)

func NewClient(opts ...ConsumerOpt) (*Client, error) {
    consumer := &Client{}
    err := applyOptionsOrDefault(consumer, opts)
    if err != nil {
        return nil, err
    }
    err = consumer.init()
    if err != nil {
        return nil, err
    }
    return consumer, nil
}

func (r *Client) Consume() (<-chan messages.Message, error) {
    if r.connChannel == nil {
        return nil, fmt.Errorf("Client was not properly initialized")
    }

    if !r.useDefaultExchange {
        if len(r.consumerRoutingKeys) == 0 {
            return nil, fmt.Errorf("missing routing key")
        }
        for _, key := range r.consumerRoutingKeys {
            err := r.connChannel.QueueBind(
                r.queue.Name,   // queue name
                key,            // routing key
                r.exchangeName, // exchange
                false,
                nil)
            if err != nil {
                fmt.Errorf("failed to bind queue %s with exchange %s using routing key %s", r.queue.Name, r.exchangeName, key)
            }
        }
    }

    msgs, err := r.connChannel.Consume(
        r.queue.Name, // queue
        "",           // consumer
        false,        // auto-ack
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
    )
    if err != nil {
        return nil, fmt.Errorf("failed to register a consumer: %w", err)
    }

    messagesCh := make(chan messages.Message)
    go func() {
        for msg := range msgs {
            messagesCh <- messages.NewMessage(msg.Body, NewAcknowledger(msg))
        }
        close(messagesCh)
    }()

    return messagesCh, nil
}

// TODO Move out of the client to a dedicated Publisher
func (r *Client) Publish(ctx context.Context, eventData events.EventData, routingInfo events.RoutingInfo) error {
    serializedEvent, err := eventData.Serialize()
    if err != nil {
        return fmt.Errorf("error serializing event: %w", err)
    }

    err = r.connChannel.PublishWithContext(ctx,
        routingInfo.Topic,      // exchange
        routingInfo.RoutingKey, // routing key
        false,                  // mandatory
        false,                  // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        serializedEvent,
        })
    if err != nil {
        return err
    }

    return nil
}

func (r *Client) Close() error {
    err := r.connChannel.Close()
    if err != nil {
        return fmt.Errorf("failed to close rabbitmq connection channel: %w", err)
    }
    err = r.conn.Close()
    if err != nil {
        return fmt.Errorf("failed to close rabbitmq connection: %w", err)
    }
    return nil
}

func (r *Client) QueueName() string {
    return r.queue.Name
}

func (r *Client) init() error {
    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", r.user, r.password, r.host, r.port))
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    r.conn = conn

    ch, err := conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open a channel: %w", err)
    }

    r.connChannel = ch

    if !r.useDefaultExchange {
        err = ch.ExchangeDeclare(
            r.exchangeName, // name
            r.exchangeType, // type
            true,           // durable
            false,          // auto-deleted
            false,          // internal
            false,          // no-wait
            nil,            // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare exchange: %w", err)
        }
    }

    q, err := ch.QueueDeclare(
        r.queueName, // name
        false,       // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare a queue: %w", err)
    }

    r.queue = q

    return nil
}

func applyOptionsOrDefault(consumer *Client, opts []ConsumerOpt) error {
    consumer.host = DefaultHost
    consumer.port = DefaultPort
    consumer.user = DefaultUser
    consumer.password = DefaultPassword
    consumer.useDefaultExchange = true
    for _, opt := range opts {
        opt(consumer)
    }
    if !consumer.useDefaultExchange {
        if consumer.exchangeType == "" {
            return fmt.Errorf("if useDefaultExchange is false exchange type can't be empty")
        }
    }
    return nil
}

func WithHost(host string) func(c *Client) {
    return func(c *Client) {
        c.host = host
    }
}

func WithPort(port uint) func(c *Client) {
    return func(c *Client) {
        c.port = port
    }
}

func WithUser(user string) func(c *Client) {
    return func(c *Client) {
        c.user = user
    }
}

func WithPassword(password string) func(c *Client) {
    return func(c *Client) {
        c.password = password
    }
}

func WithQueueName(queueName string) func(c *Client) {
    return func(c *Client) {
        c.queueName = queueName
    }
}

func WithExchangeName(exchangeName string) func(c *Client) {
    return func(c *Client) {
        c.useDefaultExchange = false
        c.exchangeName = exchangeName
    }
}

func WithExchangeType(exchangeType string) func(c *Client) {
    return func(c *Client) {
        c.exchangeType = exchangeType
    }
}

func WithConsumerRoutingKeys(routingKeys ...string) func(c *Client) {
    return func(c *Client) {
        c.consumerRoutingKeys = routingKeys
    }
}
