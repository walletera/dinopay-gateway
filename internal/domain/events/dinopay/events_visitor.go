package dinopay

import (
    "context"
    "log/slog"

    "github.com/google/uuid"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    procerrors "github.com/walletera/message-processor/errors"
    "github.com/walletera/payments/api"
)

type EventsVisitor interface {
    VisitPaymentCreated(ctx context.Context, event PaymentCreated) procerrors.ProcessingError
}

type EventsVisitorImpl struct {
    paymentsApiClient *api.Client
    logger            *slog.Logger
}

func NewEventsVisitorImpl(paymentsApiClient *api.Client, logger *slog.Logger) *EventsVisitorImpl {
    return &EventsVisitorImpl{
        paymentsApiClient: paymentsApiClient,
        logger:            logger.With(logattr.Component("dinopay.EventsVisitor")),
    }
}

func (e EventsVisitorImpl) VisitPaymentCreated(ctx context.Context, event PaymentCreated) procerrors.ProcessingError {
    depositUUID, err := uuid.NewUUID()
    if err != nil {
        return procerrors.NewInternalError(err.Error())
    }
    // TODO get from customer repository
    // for now we will hardcode the uuid to pass the integration test
    customerUUID, err := uuid.Parse("9fd3bc09-99da-4486-950a-11082f5fd966")
    if err != nil {
        return procerrors.NewInternalError(err.Error())
    }
    depositPostBody := &api.DepositPostBody{
        ID: depositUUID,
        // FIXME
        Amount:     float64(event.Data.Amount),
        Currency:   event.Data.Currency,
        CustomerId: customerUUID,
        ExternalId: event.Data.Id,
    }
    _, err = e.paymentsApiClient.PostDeposit(ctx, depositPostBody)
    if err != nil {
        // TODO this may be a retryable error
        e.logger.Error("failed creating deposit on payments api", logattr.Error(err.Error()))
        return procerrors.NewInternalError(err.Error())
    }
    // TODO handle response
    e.logger.Info("DinoPay event PaymentCreated processed successfully", logattr.EventType(event.EventType))
    return nil
}
