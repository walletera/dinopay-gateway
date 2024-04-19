package app

import (
    "context"
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/adapters/dinopay"
    "github.com/walletera/dinopay-gateway/internal/domain/events/visitors/payments"
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

    dinopayClient, err := dinopay.NewClient(a.dinopayUrl)
    if err != nil {
        return fmt.Errorf("error parsing dinopay url %s: %w", a.dinopayUrl, err)
    }

    processor, err := paymentslib.NewRabbitMQProcessor(payments.NewEventsVisitor(dinopayClient), RabbitMQQueueName)
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
