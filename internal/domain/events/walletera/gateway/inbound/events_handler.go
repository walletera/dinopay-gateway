package inbound

import (
	"context"
	"log/slog"

	"github.com/walletera/dinopay-gateway/pkg/logattr"
	procerrors "github.com/walletera/message-processor/errors"
	paymentsapi "github.com/walletera/payments/api"
)

type EventsHandler interface {
	VisitInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) procerrors.ProcessingError
}

type EventsHandlerImpl struct {
	paymentsApiClient *paymentsapi.Client
	deserializer      *EventsDeserializer
	logger            *slog.Logger
}

func NewEventsHandlerImpl(client *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
	return &EventsHandlerImpl{
		paymentsApiClient: client,
		deserializer:      NewEventsDeserializer(),
		logger:            logger.With(logattr.Component("gateway.inbound.EventsHandlerImpl")),
	}
}

func (ev *EventsHandlerImpl) VisitInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) procerrors.ProcessingError {
	depositPostBody := &paymentsapi.DepositPostBody{
		ID: inboundPaymentReceived.DepositId,
		// FIXME
		Amount:     float64(inboundPaymentReceived.Amount),
		Currency:   inboundPaymentReceived.Currency,
		CustomerId: inboundPaymentReceived.CustomerId,
		ExternalId: inboundPaymentReceived.DinopayPaymentId,
	}
	_, err := ev.paymentsApiClient.PostDeposit(ctx, depositPostBody)
	if err != nil {
		// TODO this may be a retryable error
		ev.logger.Error("failed creating deposit on payments api", logattr.Error(err.Error()))
		return procerrors.NewInternalError(err.Error())
	}
	// TODO handle response
	ev.logger.Info("Gateway event InboundPaymentReceived processed successfully", logattr.EventType(inboundPaymentReceived.Type()))
	return nil
}
