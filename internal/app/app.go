package app

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/walletera/dinopay-gateway/internal/adapters/dinopay"
    esdbAdapter "github.com/walletera/dinopay-gateway/internal/adapters/eventstoredb"
    dinopay2 "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/payments"
    esdbPkg "github.com/walletera/dinopay-gateway/pkg/eventstoredb"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    "github.com/walletera/message-processor/messages"
    paymentslib "github.com/walletera/message-processor/payments"
    paymentsApi "github.com/walletera/payments/api"
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
)

const (
    RabbitMQQueueName                         = "dinopay-gateway"
    ESDB_ByCategoryProjection_OutboundPayment = "$ce-outboundPayment"
    ESDB_SubscriptionGroupName                = "dinopay-gateway"
)

type App struct {
    rabbitmqUrl string
    dinopayUrl  string
    paymentsUrl string
    esdbUrl     string
}

func NewApp(opts ...Option) *App {
    app := &App{}
    for _, opt := range opts {
        opt(app)
    }
    return app
}

func (app *App) Run(ctx context.Context) error {
    zapLogger, err := newZapLogger()
    if err != nil {
        return err
    }

    appLogger := slog.
        New(zapslog.NewHandler(zapLogger.Core(), nil)).
        With(logattr.ServiceName("dinopay-gateway"))

    paymentsMessageProcessor, err := createPaymentsMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating payments message processor: %w", err)
    }

    err = paymentsMessageProcessor.Start()
    if err != nil {
        return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
    }

    dinopayMessageProcessor, err := createDinopayMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating dinopay message processor: %w", err)
    }

    err = dinopayMessageProcessor.Start()
    if err != nil {
        return fmt.Errorf("failed starting dinopay message processor: %w", err)
    }

    appLogger.Info("dinopay-gateway started")
    <-ctx.Done()

    appLogger.Info("dinopay-gateway stopped")
    return nil
}

func newZapLogger() (*zap.Logger, error) {
    zapConfig := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:    zap.NewProductionEncoderConfig(),
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
    return zapConfig.Build()
}

func createPaymentsMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[paymentslib.EventsVisitor], error) {
    dinopayClient, err := dinopay.NewClient(app.dinopayUrl)
    if err != nil {
        return nil, fmt.Errorf("failed parsing dinopay url %s: %w", app.dinopayUrl, err)
    }

    esdbClient, err := esdbPkg.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed getting esdb client: %w", err)
    }

    eventsDB := esdbAdapter.NewDB(esdbClient)
    visitor := payments.NewEventsVisitor(dinopayClient, eventsDB, logger)
    paymentsMessageProcessor, err := paymentslib.NewRabbitMQProcessor(visitor, RabbitMQQueueName)
    if err != nil {
        return nil, fmt.Errorf("failed creating payments rabbitmq processor: %w", err)
    }

    return paymentsMessageProcessor, nil
}

func createDinopayMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[dinopay2.EventsVisitor], error) {

    paymentsClient, err := paymentsApi.NewClient(app.paymentsUrl)
    if err != nil {
        return nil, fmt.Errorf("failed creating payments api client: %w", err)
    }

    esdbMessagesConsumer, err := esdbPkg.NewMessagesConsumer(
        app.esdbUrl,
        ESDB_ByCategoryProjection_OutboundPayment,
        ESDB_SubscriptionGroupName,
    )
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb messages consumer: %w", err)
    }

    esdbClient, err := esdbPkg.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb client: %w", err)
    }

    eventsDB := esdbAdapter.NewDB(esdbClient)

    eventsVisitor := dinopay2.NewEventsVisitorImpl(eventsDB, paymentsClient, logger)
    return messages.NewProcessor[dinopay2.EventsVisitor](
        esdbMessagesConsumer,
        dinopay2.NewEventsDeserializer(),
        eventsVisitor,
    ), nil
}
