package outbound

import (
    "context"
    "errors"
    "log/slog"

    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/inbound"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    "github.com/walletera/eventskit/eventsourcing"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

type EventsHandler interface {
    HandleOutboundPaymentCreated(ctx context.Context, outboundPaymentCreated PaymentCreated) werrors.WError
    HandleOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated PaymentUpdated) werrors.WError
}

type EventsHandlerImpl struct {
    db             eventsourcing.DB
    paymentsClient *paymentsapi.Client
    deserializer   *EventsDeserializer
    logger         *slog.Logger
}

func NewEventsHandlerImpl(db eventsourcing.DB, client *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
        logger:         logger.With(logattr.Component("dinopay.EventsHandlerImpl")),
    }
}

func (ev *EventsHandlerImpl) HandleOutboundPaymentCreated(ctx context.Context, outboundPaymentCreated PaymentCreated) werrors.WError {
    err := NewOutboundPaymentCreatedHandler(ev.paymentsClient).Handle(ctx, outboundPaymentCreated)
    if err != nil {
        ev.logger.Error(
            err.Error(),
            logattr.EventType(outboundPaymentCreated.Type()),
            logattr.PaymentId(outboundPaymentCreated.PaymentId.String()))
        return werrors.NewWrappedError(err, "failed handling outbound PaymentCreated event")
    }
    ev.logger.Info("OutboundPaymentCreated event processed successfully", logattr.PaymentId(outboundPaymentCreated.PaymentId.String()))
    return nil
}

func (ev *EventsHandlerImpl) HandleOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated PaymentUpdated) werrors.WError {
    err := NewOutboundPaymentUpdatedHandler(ev.db, ev.paymentsClient).Handle(ctx, outboundPaymentUpdated)
    if err != nil {
        logOutboundPaymentUpdatedHandlerError(ev.logger, outboundPaymentUpdated, err)
        return werrors.NewWrappedError(err, "failed handling outbound PaymentUpdated event")
    }
    ev.logger.Info("OutboundPaymentUpdated event processed successfully", logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()))
    return nil
}

func (ev *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived inbound.PaymentReceived) werrors.WError {
    //err := NewInboundPaymentReceivedHandler(ev.db, ev.paymentsClient).Handle(ctx, inboundPaymentReceived)
    //if err != nil {
    //    logOutboundPaymentUpdatedHandlerError(ev.logger, inboundPaymentReceived, err)
    //    return procerrors.NewInternalError(err.Error())
    //}
    //ev.logger.Info("OutboundPaymentUpdated event processed successfully", logattr.DinopayPaymentId(inboundPaymentReceived.DinopayPaymentId.String()))
    return nil
}

func logOutboundPaymentUpdatedHandlerError(logger *slog.Logger, outboundPaymentUpdated PaymentUpdated, err error) {
    logger = logger.With(
        logattr.EventType(outboundPaymentUpdated.Type()),
        logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()),
    )
    var handlerErr *HandlerError
    isHandlerErr := errors.As(err, &handlerErr)
    if isHandlerErr {
        logger = logger.With(logattr.PaymentId(handlerErr.withdrawalId))
    }
    logger.Error(err.Error())
}
