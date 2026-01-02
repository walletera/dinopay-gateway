package dinopay

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	accountsapi "github.com/walletera/accounts/publicapi"
	"github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/inbound"
	gatewayevents "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/outbound"
	"github.com/walletera/dinopay-gateway/pkg/logattr"
	"github.com/walletera/dinopay-gateway/pkg/wuuid"
	"github.com/walletera/eventskit/eventsourcing"
	paymentsapi "github.com/walletera/payments-types/privateapi"
	"github.com/walletera/werrors"
)

type EventsHandler interface {
	HandlePaymentCreated(ctx context.Context, event PaymentCreated) werrors.WError
}

type EventsHandlerImpl struct {
	db                eventsourcing.DB
	accountsApiClient *accountsapi.Client
	paymentsApiClient *paymentsapi.Client
	logger            *slog.Logger
}

func NewEventsHandlerImpl(
	db eventsourcing.DB,
	accountsApiClient *accountsapi.Client,
	paymentsApiClient *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
	return &EventsHandlerImpl{
		db:                db,
		accountsApiClient: accountsApiClient,
		paymentsApiClient: paymentsApiClient,
		logger:            logger.With(logattr.Component("dinopay.EventsHandler")),
	}
}

func (ev EventsHandlerImpl) HandlePaymentCreated(ctx context.Context, event PaymentCreated) werrors.WError {
	// TODO get correlation id from PaymentCreated and copy into PaymentReceived
	eventUUID := wuuid.NewUUID()
	depositUUID := wuuid.NewUUID()
	accountNumber := event.Data.DestinationAccount.AccountNumber
	resp, err := ev.accountsApiClient.ListAccounts(ctx, accountsapi.ListAccountsParams{DinopayAccountNumber: accountsapi.NewOptString(accountNumber)})
	if err != nil {
		ev.logger.Error("failed to list accounts", logattr.Error(err.Error()))
		return werrors.NewRetryableInternalError("failed to list accounts for dinopay payment: %s", err.Error())
	}
	var customerUUID uuid.UUID
	switch resp.(type) {
	case *accountsapi.ListAccountsOKApplicationJSON:
		accountList := resp.(*accountsapi.ListAccountsOKApplicationJSON)
		if len(*accountList) == 0 {
			ev.logger.With("account-number", accountNumber).Error("no account found")
			return werrors.NewNonRetryableInternalError("no account found")
		}
		if len(*accountList) > 1 {
			ev.logger.With("account-number", accountNumber).Error("multiple accounts found")
			return werrors.NewNonRetryableInternalError("multiple accounts found")
		}
		customerUUID = (*accountList)[0].CustomerId
	case *accountsapi.ListAccountsUnauthorized:
		ev.logger.Error("unauthorized to list accounts")
		return werrors.NewNonRetryableInternalError("unauthorized to list accounts")
	case *accountsapi.ApiError:
		ev.logger.Error("error listing accounts")
		return werrors.NewNonRetryableInternalError("error listing accounts")
	default:
		ev.logger.Error("unknown response type")
		return werrors.NewNonRetryableInternalError("unknown response type")
	}

	inboundPaymentReceived := inbound.PaymentReceived{
		Id:               eventUUID,
		DinopayPaymentId: event.Data.Id,
		CustomerId:       customerUUID,
		PaymentId:        depositUUID,
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
		EventCreatedAt: time.Now(),
	}
	streamName := gatewayevents.BuildInboundPaymentStreamName(inboundPaymentReceived.DinopayPaymentId.String())
	_, werr := ev.db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, inboundPaymentReceived)
	if werr != nil {
		ev.logger.Error("error handling dinopay PaymentCreated event", logattr.Error(werr.Error()))
		return werrors.NewWrappedError(werr)
	}
	ev.logger.Info("DinoPay event PaymentCreated processed successfully", logattr.EventType(event.Type()))
	return nil
}
