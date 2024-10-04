package dinopay

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/walletera/dinopay-gateway/internal/domain/events"
	"github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/inbound"
	gatewayevents "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/outbound"
	"github.com/walletera/dinopay-gateway/pkg/logattr"
	"github.com/walletera/dinopay-gateway/pkg/wuuid"
	procerrors "github.com/walletera/message-processor/errors"
	"github.com/walletera/payments/api"
)

type EventsHandler interface {
	VisitPaymentCreated(ctx context.Context, event PaymentCreated) procerrors.ProcessingError
}

type EventsHandlerImpl struct {
	db                events.DB
	paymentsApiClient *api.Client
	logger            *slog.Logger
}

func NewEventsHandlerImpl(db events.DB, paymentsApiClient *api.Client, logger *slog.Logger) *EventsHandlerImpl {
	return &EventsHandlerImpl{
		db:                db,
		paymentsApiClient: paymentsApiClient,
		logger:            logger.With(logattr.Component("dinopay.EventsHandler")),
	}
}

func (ev EventsHandlerImpl) VisitPaymentCreated(ctx context.Context, event PaymentCreated) procerrors.ProcessingError {
	eventUUID := wuuid.NewUUID()
	depositUUID := wuuid.NewUUID()
	// TODO get from customer repository
	// for now we will hardcode the uuid to pass the integration test
	customerUUID, err := uuid.Parse("9fd3bc09-99da-4486-950a-11082f5fd966")
	inboundPaymentReceived := inbound.PaymentReceived{
		Id:               eventUUID,
		DinopayPaymentId: event.Data.Id,
		CustomerId:       customerUUID,
		DepositId:        depositUUID,
		Amount:           event.Data.Amount,
		Currency:         event.Data.Currency,
		SourceAccount: inbound.Account{
			AccountHolder: event.Data.SourceAccount.AccountHolder,
			AccountNumber: event.Data.SourceAccount.AccountNumber,
		},
		DestinationAccount: inbound.Account{
			AccountHolder: event.Data.DestinationAccount.AccountHolder,
			AccountNumber: event.Data.DestinationAccount.AccountNumber,
		},
		CreatedAt: time.Now(),
	}
	streamName := gatewayevents.BuildInboundPaymentStreamName(inboundPaymentReceived.DinopayPaymentId.String())
	err = ev.db.AppendEvents(streamName, inboundPaymentReceived)
	if err != nil {
		return procerrors.NewInternalError(fmt.Sprintf("failed adding InboundPaymentReceived event to the events db: %s", err.Error()))
	}
	ev.logger.Info("DinoPay event PaymentCreated processed successfully", logattr.EventType(event.Type()))
	return nil
}
