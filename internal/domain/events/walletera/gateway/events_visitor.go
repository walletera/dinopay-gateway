package gateway

import (
    "errors"
    "log/slog"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    "github.com/walletera/dinopay-gateway/pkg/logattr"
    paymentsApi "github.com/walletera/payments/api"
)

type EventsVisitor interface {
    VisitOutboundPaymentCreated(outboundPaymentCreated OutboundPaymentCreated) error
    VisitOutboundPaymentUpdated(outboundPaymentUpdated OutboundPaymentUpdated) error
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
        logger:         logger.With(logattr.Component("dinopay/EventsVisitorImpl")),
    }
}

func (ev *EventsVisitorImpl) VisitOutboundPaymentCreated(outboundPaymentCreated OutboundPaymentCreated) error {
    err := NewOutboundPaymentCreatedHandler(ev.paymentsClient).Handle(outboundPaymentCreated)
    if err != nil {
        ev.logger.Error(
            err.Error(),
            logattr.EventType(outboundPaymentCreated.Type()),
            logattr.WithdrawalId(outboundPaymentCreated.WithdrawalId.String()))
        return err
    }
    ev.logger.Info("OutboundPaymentCreated event processed successfully", logattr.WithdrawalId(outboundPaymentCreated.WithdrawalId.String()))
    return nil
}

func (ev *EventsVisitorImpl) VisitOutboundPaymentUpdated(outboundPaymentUpdated OutboundPaymentUpdated) error {
    err := NewOutboundPaymentUpdatedHandler(ev.db, ev.paymentsClient).Handle(outboundPaymentUpdated)
    if err != nil {
        logOutboundPaymentUpdatedHandlerError(ev.logger, outboundPaymentUpdated, err)
        return err
    }
    ev.logger.Info("OutboundPaymentUpdated event processed successfully", logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()))
    return nil
}

func logOutboundPaymentUpdatedHandlerError(logger *slog.Logger, outboundPaymentUpdated OutboundPaymentUpdated, err error) {
    loggerAttrs := []slog.Attr{
        logattr.EventType(outboundPaymentUpdated.Type()),
        logattr.DinopayPaymentId(outboundPaymentUpdated.DinopayPaymentId.String()),
    }
    var handlerErr *HandlerError
    ok := errors.As(err, &handlerErr)
    if ok {
        loggerAttrs = append(loggerAttrs, logattr.WithdrawalId(handlerErr.withdrawalId))
    }
    logger.Error(err.Error(), loggerAttrs)
}
