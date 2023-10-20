package app

import (
    "context"
    "fmt"
    "github.com/walletera/dinopay-gateway/internal/messages/processors/payments"
    paymentslib "github.com/walletera/message-processor/pkg/events/payments"
    "go.uber.org/zap"
)

const (
    RabbitMQQueueName = "dinopay-gateway"
)

type App struct {
    rabbitmqUrl string
    dinopayUrl  string
    paymentsUrl string
}

func NewApp(opts ...Option) *App {
    app := &App{}
    for _, opt := range opts {
        opt(app)
    }
    return app
}

func (a *App) Run(ctx context.Context) error {
    logger, err := zap.NewDevelopment()
    if err != nil {
        return fmt.Errorf("error creating logger: %w", err)
    }

    processor, err := paymentslib.NewRabbitMQProcessor(payments.NewEventsVisitor(), RabbitMQQueueName)
    if err != nil {
        return fmt.Errorf("error creating payments rabbitmq processor: %w", err)
    }

    err = processor.Start()
    if err != nil {
        return fmt.Errorf("error starting payments rabbitmq processor: %w", err)
    }

    logger.Info("dinopay-gateway started")
    <-ctx.Done()

    logger.Info("dinopay-gateway stopped")
    return nil
}
