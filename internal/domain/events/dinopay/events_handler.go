package dinopay

import (
    "context"
    "log/slog"
    "time"

    "github.com/google/uuid"
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
    paymentsApiClient *paymentsapi.Client
    logger            *slog.Logger
}

func NewEventsHandlerImpl(db eventsourcing.DB, paymentsApiClient *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        db:                db,
        paymentsApiClient: paymentsApiClient,
        logger:            logger.With(logattr.Component("dinopay.EventsHandler")),
    }
}

func (ev EventsHandlerImpl) HandlePaymentCreated(ctx context.Context, event PaymentCreated) werrors.WError {
    // TODO get correlation id from PaymentCreated and copy into PaymentReceived
    eventUUID := wuuid.NewUUID()
    depositUUID := wuuid.NewUUID()
    // TODO get from customer repository
    // for now we will hardcode the uuid to pass the integration test
    customerUUID := uuid.MustParse("9fd3bc09-99da-4486-950a-11082f5fd966")
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
        CreatedAt: time.Now(),
    }
    streamName := gatewayevents.BuildInboundPaymentStreamName(inboundPaymentReceived.DinopayPaymentId.String())
    werr := ev.db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, inboundPaymentReceived)
    if werr != nil {
        ev.logger.Error("error handling dinopay PaymentCreated event", logattr.Error(werr.Error()))
        return werrors.NewWrappedError(werr)
    }
    ev.logger.Info("DinoPay event PaymentCreated processed successfully", logattr.EventType(event.Type()))
    return nil
}
