package app

import (
    "context"
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/adapters/dinopay"
    esdbAdapter "github.com/walletera/dinopay-gateway/internal/adapters/eventstoredb"
    dinopayEvents "github.com/walletera/dinopay-gateway/internal/domain/events/dinopay"
    "github.com/walletera/dinopay-gateway/internal/domain/events/payments"
    esdbPkg "github.com/walletera/dinopay-gateway/pkg/eventstoredb"
    "github.com/walletera/message-processor/messages"
    paymentslib "github.com/walletera/message-processor/payments"
    paymentsApi "github.com/walletera/payments/api"
    "go.uber.org/zap"
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
    logger, err := zap.NewDevelopment()
    if err != nil {
        return fmt.Errorf("error creating logger: %w", err)
    }

    paymentsMessageProcessor, err := createPaymentsMessageProcessor(app)
    if err != nil {
        return fmt.Errorf("failed creating payments message processor: %w", err)
    }

    err = paymentsMessageProcessor.Start()
    if err != nil {
        return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
    }

    dinopayMessageProcessor, err := createDinopayMessageProcessor(app)
    if err != nil {
        return fmt.Errorf("failed creating dinopay message processor: %w", err)
    }

    err = dinopayMessageProcessor.Start()
    if err != nil {
        return fmt.Errorf("failed starting dinopay message processor: %w", err)
    }

    logger.Info("dinopay-gateway started")
    <-ctx.Done()

    logger.Info("dinopay-gateway stopped")
    return nil
}

func createPaymentsMessageProcessor(app *App) (*messages.Processor[paymentslib.EventsVisitor], error) {
    dinopayClient, err := dinopay.NewClient(app.dinopayUrl)
    if err != nil {
        return nil, fmt.Errorf("failed parsing dinopay url %s: %w", app.dinopayUrl, err)
    }

    esdbClient, err := esdbPkg.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed getting esdb client: %w", err)
    }

    eventsDB := esdbAdapter.NewDB(esdbClient)

    paymentsMessageProcessor, err := paymentslib.NewRabbitMQProcessor(payments.NewEventsVisitor(dinopayClient, eventsDB), RabbitMQQueueName)
    if err != nil {
        return nil, fmt.Errorf("failed creating payments rabbitmq processor: %w", err)
    }

    return paymentsMessageProcessor, nil
}

func createDinopayMessageProcessor(app *App) (*messages.Processor[dinopayEvents.EventsVisitor], error) {

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

    return messages.NewProcessor[dinopayEvents.EventsVisitor](
        esdbMessagesConsumer,
        dinopayEvents.NewEventsDeserializer(),
        dinopayEvents.NewEventsVisitorImpl(eventsDB, paymentsClient),
    ), nil
}
