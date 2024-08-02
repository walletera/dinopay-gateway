package outbound

import (
    "context"
    "errors"
    "log/slog"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    "github.com/walletera/dinopay-gateway/internal/domain/events/walletera/gateway/inbound"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    procerrors "github.com/walletera/message-processor/errors"
    paymentsApi "github.com/walletera/payments/api"
)

type EventsVisitor interface {
    VisitOutboundPaymentCreated(ctx context.Context, outboundPaymentCreated OutboundPaymentCreated) procerrors.ProcessingError
    VisitOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated OutboundPaymentUpdated) procerrors.ProcessingError
}

type EventsVisitorImpl struct {
    db             events.DB
    paymentsClient *paymentsApi.Client
    deserializer   *EventsDeserializer
    logger         *slog.Logger
}

func NewEventsVisitorImpl(db events.DB, client *paymentsApi.Client, logger *slog.Logger) *EventsVisitorImpl {
    return &EventsVisitorImpl{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
        logger:         logger.With(logattr.Component("dinopay.EventsVisitorImpl")),
    }
}

func (ev *EventsVisitorImpl) VisitOutboundPaymentCreated(ctx context.Context, outboundPaymentCreated OutboundPaymentCreated) procerrors.ProcessingError {
    err := NewOutboundPaymentCreatedHandler(ev.paymentsClient).Handle(ctx, outboundPaymentCreated)
    if err != nil {
        ev.logger.Error(
            err.Error(),
            logattr.EventType(outboundPaymentCreated.Type()),
            logattr.WithdrawalId(outboundPaymentCreated.WithdrawalId.String()))
        return procerrors.NewInternalError(err.Error())
    }
    ev.logger.Info("OutboundPaymentCreated event processed successfully", logattr.WithdrawalId(outboundPaymentCreated.WithdrawalId.String()))
    return nil
}

func (ev *EventsVisitorImpl) VisitOutboundPaymentUpdated(ctx context.Context, outboundPaymentUpdated OutboundPaymentUpdated) procerrors.ProcessingError {
    err := NewOutboundPaymentUpdatedHandler(ev.db, ev.paymentsClient).Handle(ctx, outboundPaymentUpdated)
    if err != nil {
        logOutboundPaymentUpdatedHandlerError(ev.logger, outboundPaymentUpdated, err)
        return procerrors.NewInternalError(err.Error())
    }
    ev.logger.Info("OutboundPaymentUpdated event processed successfully", logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()))
    return nil
}

func (ev *EventsVisitorImpl) VisitInboundPaymentReceived(ctx context.Context, inboundPaymentReceived inbound.PaymentReceived) procerrors.ProcessingError {
    //err := NewInboundPaymentReceivedHandler(ev.db, ev.paymentsClient).Handle(ctx, inboundPaymentReceived)
    //if err != nil {
    //    logOutboundPaymentUpdatedHandlerError(ev.logger, inboundPaymentReceived, err)
    //    return procerrors.NewInternalError(err.Error())
    //}
    //ev.logger.Info("OutboundPaymentUpdated event processed successfully", logattr.DinopayPaymentId(inboundPaymentReceived.DinopayPaymentId.String()))
    return nil
}

func logOutboundPaymentUpdatedHandlerError(logger *slog.Logger, outboundPaymentUpdated OutboundPaymentUpdated, err error) {
    logger = logger.With(
        logattr.EventType(outboundPaymentUpdated.Type()),
        logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()),
    )
    var handlerErr *HandlerError
    isHandlerErr := errors.As(err, &handlerErr)
    if isHandlerErr {
        logger = logger.With(logattr.WithdrawalId(handlerErr.withdrawalId))
    }
    logger.Error(err.Error())
}
