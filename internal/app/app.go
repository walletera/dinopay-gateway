package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EventStore/EventStore-Client-Go/v4/esdb"
	accountsapi "github.com/walletera/accounts/publicapi"
	"github.com/walletera/dinopay-gateway/internal/adapters/dinopay"
	dinopayevents "github.com/walletera/dinopay-gateway/internal/domain/events/dinopay"
	"github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/inbound"
	"github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/outbound"
	"github.com/walletera/dinopay-gateway/internal/domain/events/walletera/payments"
	"github.com/walletera/dinopay-gateway/pkg/logattr"
	"github.com/walletera/eventskit/eventstoredb"
	"github.com/walletera/eventskit/messages"
	"github.com/walletera/eventskit/rabbitmq"
	"github.com/walletera/eventskit/webhook"
	paymentsevents "github.com/walletera/payments-types/events"
	paymentsapi "github.com/walletera/payments-types/privateapi"
	"github.com/walletera/werrors"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

const (
	RabbitMQPaymentsExchangeName              = "payments.events"
	RabbitMQExchangeType                      = "topic"
	RabbitMQPaymentCreatedRoutingKey          = "payment.created"
	RabbitMQQueueName                         = "dinopay-gateway"
	ESDB_ByCategoryProjection_OutboundPayment = "$ce-outboundPayment"
	ESDB_ByCategoryProjection_InboundPayment  = "$ce-inboundPayment"
	ESDB_SubscriptionGroupName                = "dinopay-gateway"
	WebhookServerPort                         = 8686
)

type App struct {
	rabbitmqHost     string
	rabbitmqPort     int
	rabbitmqUser     string
	rabbitmqPassword string
	dinopayUrl       string
	accountsUrl      string
	paymentsUrl      string
	esdbUrl          string
	webhookConsumer  *webhook.Server
	esdbClient       *esdb.Client
	rabbitMQClient   *rabbitmq.Client
	logHandler       slog.Handler
	logger           *slog.Logger
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
	app.logger = appLogger

	err := app.execESDBSetupTasks(ctx)
	if err != nil {
		return err
	}

	paymentsMessageProcessor, err := createPaymentsMessageProcessor(app, appLogger)
	if err != nil {
		return fmt.Errorf("failed creating payments message processor: %w", err)
	}

	err = paymentsMessageProcessor.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
	}

	appLogger.Info("payments message processor started")

	dinopayMessageProcessor, err := createDinopayMessageProcessor(app, appLogger)
	if err != nil {
		return fmt.Errorf("failed creating dinopay webhook message processor: %w", err)
	}

	err = dinopayMessageProcessor.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
	}

	appLogger.Info("dinopay message processor started")

	gatewayMessageProcessor, err := createGatewayMessageProcessor(app, appLogger)
	if err != nil {
		return fmt.Errorf("failed creating dinopay message processor: %w", err)
	}

	err = gatewayMessageProcessor.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed starting dinopay message processor: %w", err)
	}

	inboundMessageProcessor, err := createGatewayInboundMessageProcessor(app, appLogger)
	if err != nil {
		return fmt.Errorf("failed creating gateway inbound message processor: %w", err)
	}

	err = inboundMessageProcessor.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed starting gateway inbound message processor: %w", err)
	}

	appLogger.Info("gateway message processor started")

	appLogger.Info("dinopay-gateway started")

	return nil
}

func (app *App) Stop(ctx context.Context) {
	if app.webhookConsumer != nil {
		err := app.webhookConsumer.Close()
		if err != nil {
			app.logger.With(logattr.Error(err.Error())).Error("failed closing webhook consumer")
		}
	}
	if app.esdbClient != nil {
		err := app.esdbClient.Close()
		if err != nil {
			app.logger.With(logattr.Error(err.Error())).Error("failed closing esdb client")
		}
	}
	if app.rabbitMQClient != nil {
		err := app.rabbitMQClient.Close()
		if err != nil {
			app.logger.With(logattr.Error(err.Error())).Error("failed closing rabbitmq client")
		}
	}
	app.logger.Info("dinopay-gateway stopped")
}

func setDefaultOpts(app *App) error {
<<<<<<< Updated upstream
    zapLogger, err := newZapLogger()
    if err != nil {
        return err
    }
    app.logHandler = zapslog.NewHandler(
        zapLogger.Core(),
        // never add stacktrace
        zapslog.AddStacktraceAt(slog.LevelError+1),
    )
    return nil
=======
	zapLogger, err := newZapLogger()
	if err != nil {
		return err
	}
	app.logHandler = zapslog.NewHandler(
		zapLogger.Core(),
		// never add stacktrace
		zapslog.AddStacktraceAt(slog.LevelError+1),
	)
	return nil
>>>>>>> Stashed changes
}

func newZapLogger() (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       false,
		DisableStacktrace: true,
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

func (app *App) execESDBSetupTasks(_ context.Context) error {
	subscriptionSettings := esdb.SubscriptionSettingsDefault()
	subscriptionSettings.ResolveLinkTos = true

	err := eventstoredb.CreatePersistentSubscription(
		app.esdbUrl,
		ESDB_ByCategoryProjection_OutboundPayment,
		ESDB_SubscriptionGroupName,
		subscriptionSettings,
	)
	if err != nil {
		return fmt.Errorf("failed creating persistent subscription for %s: %w", ESDB_ByCategoryProjection_OutboundPayment, err)
	}

	err = eventstoredb.CreatePersistentSubscription(
		app.esdbUrl,
		ESDB_ByCategoryProjection_InboundPayment,
		ESDB_SubscriptionGroupName,
		subscriptionSettings,
	)
	if err != nil {
		return fmt.Errorf("failed creating persistent subscription for %s: %w", ESDB_ByCategoryProjection_InboundPayment, err)
	}
	return nil
}

func createPaymentsMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[paymentsevents.Handler], error) {
	dinopayClient, err := dinopay.NewClient(app.dinopayUrl)
	if err != nil {
		return nil, fmt.Errorf("failed parsing dinopay url %s: %w", app.dinopayUrl, err)
	}

	esdbClient, err := eventstoredb.GetESDBClient(app.esdbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed getting esdb client: %w", err)
	}

	eventsDB := eventstoredb.NewDB(esdbClient)
	handler := payments.NewEventsHandler(dinopayClient, eventsDB, logger)
	queueName := fmt.Sprintf(RabbitMQQueueName)

	rabbitMQClient, err := rabbitmq.NewClient(
		rabbitmq.WithHost(app.rabbitmqHost),
		rabbitmq.WithPort(uint(app.rabbitmqPort)),
		rabbitmq.WithUser(app.rabbitmqUser),
		rabbitmq.WithPassword(app.rabbitmqPassword),
		rabbitmq.WithExchangeName(RabbitMQPaymentsExchangeName),
		rabbitmq.WithExchangeType(RabbitMQExchangeType),
		rabbitmq.WithConsumerRoutingKeys(RabbitMQPaymentCreatedRoutingKey),
		rabbitmq.WithQueueName(queueName),
	)
	if err != nil {
		return nil, fmt.Errorf("creating rabbitmq client: %w", err)
	}

	paymentsMessageProcessor := messages.NewProcessor[paymentsevents.Handler](
		rabbitMQClient,
		paymentsevents.NewDeserializer(logger),
		handler,
		withErrorCallback(
			logger.With(
				logattr.Component("payments.rabbitmq.MessageProcessor")),
		),
	)

	app.esdbClient = esdbClient
	app.rabbitMQClient = rabbitMQClient

	return paymentsMessageProcessor, nil
}

type AccountsSecuritySource struct {
}

func (a AccountsSecuritySource) BearerAuth(ctx context.Context, operationName accountsapi.OperationName) (accountsapi.BearerAuth, error) {
	//TODO implement me
	return accountsapi.BearerAuth{
		Token: "somejsonwebtoken",
		Roles: nil,
	}, nil
}

func createDinopayMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[dinopayevents.EventsHandler], error) {

	accountsapiClient, err := accountsapi.NewClient(app.accountsUrl, AccountsSecuritySource{})
	if err != nil {
		return nil, fmt.Errorf("failed creating accounts api client: %w", err)
	}
	paymentsClient, err := paymentsapi.NewClient(app.paymentsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed creating payments api client: %w", err)
	}
	webhookConsumer := webhook.NewServer(WebhookServerPort, webhook.WithLogger(logger.With(logattr.Component("webhook.Server"))))
	app.webhookConsumer = webhookConsumer

	esdbClient, err := eventstoredb.GetESDBClient(app.esdbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed getting esdb client: %w", err)
	}
	eventsDB := eventstoredb.NewDB(esdbClient)
	eventsHandler := dinopayevents.NewEventsHandlerImpl(eventsDB, accountsapiClient, paymentsClient, logger)
	return messages.NewProcessor[dinopayevents.EventsHandler](
		webhookConsumer,
		dinopayevents.NewEventsDeserializer(),
		eventsHandler,
		withErrorCallback(
			logger.With(
				logattr.Component("dinopay.webhook.MessageProcessor"),
			),
		),
	), nil
}

func createGatewayInboundMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[inbound.EventsHandler], error) {

	paymentsClient, err := paymentsapi.NewClient(app.paymentsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed creating payments api client: %w", err)
	}

	esdbMessagesConsumer, err := eventstoredb.NewMessagesConsumer(
		app.esdbUrl,
		ESDB_ByCategoryProjection_InboundPayment,
		ESDB_SubscriptionGroupName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating esdb messages consumer: %w", err)
	}

	eventsHandler := inbound.NewEventsHandlerImpl(paymentsClient, logger)
	return messages.NewProcessor[inbound.EventsHandler](
		esdbMessagesConsumer,
		inbound.NewEventsDeserializer(),
		eventsHandler,
		withErrorCallback(
			logger.With(
				logattr.Component("gateway.inbound.MessageProcessor"),
			),
		),
	), nil
}

func createGatewayMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[outbound.EventsHandler], error) {

	paymentsClient, err := paymentsapi.NewClient(app.paymentsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed creating payments api client: %w", err)
	}

	esdbMessagesConsumer, err := eventstoredb.NewMessagesConsumer(
		app.esdbUrl,
		ESDB_ByCategoryProjection_OutboundPayment,
		ESDB_SubscriptionGroupName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating esdb messages consumer: %w", err)
	}

	esdbClient, err := eventstoredb.GetESDBClient(app.esdbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed creating esdb client: %w", err)
	}

	eventsDB := eventstoredb.NewDB(esdbClient)

	eventsHandler := outbound.NewEventsHandlerImpl(eventsDB, paymentsClient, logger)
	return messages.NewProcessor[outbound.EventsHandler](
			esdbMessagesConsumer,
			outbound.NewEventsDeserializer(),
			eventsHandler,
			withErrorCallback(
				logger.With(
					logattr.Component("gateway.esdb.MessageProcessor")),
			),
		),
		nil
}

func withErrorCallback(logger *slog.Logger) messages.ProcessorOpt {
	return messages.WithErrorCallback(func(wError werrors.WError) {
		logger.Error(
			"failed processing message",
			logattr.Error(wError.Message()))
	})
}
