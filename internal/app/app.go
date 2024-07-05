package app

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/walletera/dinopay-gateway/internal/adapters/dinopay"
    esdbadapter "github.com/walletera/dinopay-gateway/internal/adapters/eventstoredb"
    gatewayevents "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/payments"
    esdbpkg "github.com/walletera/dinopay-gateway/pkg/eventstoredb"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    "github.com/walletera/message-processor/messages"
    paymentsevents "github.com/walletera/message-processor/payments"
    paymentsapi "github.com/walletera/payments/api"
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
    "go.uber.org/zap/zapcore"
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
    logHandler  slog.Handler
}

func NewApp(opts ...Option) (*App, error) {
    app := &App{}
    err := setDefaultOpts(app)
    if err != nil {
        return nil, fmt.Errorf("failed setting default options: %w", err)
    }
    for _, opt := range opts {
        opt(app)
    }
    return app, nil
}

func (app *App) Run(ctx context.Context) error {
    appLogger := slog.
        New(app.logHandler).
        With(logattr.ServiceName("dinopay-gateway"))

    paymentsMessageProcessor, err := createPaymentsMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating payments message processor: %w", err)
    }

    err = paymentsMessageProcessor.Start(ctx)
    if err != nil {
        return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
    }

    gatewayMessageProcessor, err := createGatewayMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating dinopay message processor: %w", err)
    }

    err = gatewayMessageProcessor.Start(ctx)
    if err != nil {
        return fmt.Errorf("failed starting dinopay message processor: %w", err)
    }

    appLogger.Info("dinopay-gateway started")
    <-ctx.Done()

    appLogger.Info("dinopay-gateway stopped")
    return nil
}

func setDefaultOpts(app *App) error {
    zapLogger, err := newZapLogger()
    if err != nil {
        return err
    }
    app.logHandler = zapslog.NewHandler(zapLogger.Core(), nil)
    return nil
}

func newZapLogger() (*zap.Logger, error) {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
    zapConfig := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:    encoderConfig,
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
    return zapConfig.Build()
}

func createPaymentsMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[paymentsevents.EventsVisitor], error) {
    dinopayClient, err := dinopay.NewClient(app.dinopayUrl)
    if err != nil {
        return nil, fmt.Errorf("failed parsing dinopay url %s: %w", app.dinopayUrl, err)
    }

    esdbClient, err := esdbpkg.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed getting esdb client: %w", err)
    }

    eventsDB := esdbadapter.NewDB(esdbClient)
    visitor := payments.NewEventsVisitor(dinopayClient, eventsDB, logger)
    queueName := fmt.Sprintf(RabbitMQQueueName)
    paymentsMessageProcessor, err := paymentsevents.NewRabbitMQProcessor(visitor, queueName)
    if err != nil {
        return nil, fmt.Errorf("failed creating payments rabbitmq processor: %w", err)
    }

    return paymentsMessageProcessor, nil
}

func createGatewayMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[gatewayevents.EventsVisitor], error) {

    paymentsClient, err := paymentsapi.NewClient(app.paymentsUrl)
    if err != nil {
        return nil, fmt.Errorf("failed creating payments api client: %w", err)
    }

    esdbMessagesConsumer, err := esdbpkg.NewMessagesConsumer(
        app.esdbUrl,
        ESDB_ByCategoryProjection_OutboundPayment,
        ESDB_SubscriptionGroupName,
    )
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb messages consumer: %w", err)
    }

    esdbClient, err := esdbpkg.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb client: %w", err)
    }

    eventsDB := esdbadapter.NewDB(esdbClient)

    eventsVisitor := gatewayevents.NewEventsVisitorImpl(eventsDB, paymentsClient, logger)
    return messages.NewProcessor[gatewayevents.EventsVisitor](
        esdbMessagesConsumer,
        gatewayevents.NewEventsDeserializer(),
        eventsVisitor,
    ), nil
}
