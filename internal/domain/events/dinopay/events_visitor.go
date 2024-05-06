package dinopay

import (
    "fmt"

    "github.com/walletera/dinopay-gateway/internal/domain/events"
    paymentsApi "github.com/walletera/payments/api"
    "go.uber.org/zap"
)

type EventsVisitor interface {
    VisitOutboundPaymentCreated(outboundPaymentCreated OutboundPaymentCreated) error
    VisitOutboundPaymentUpdated(outboundPaymentUpdated OutboundPaymentUpdated) error
}

type EventsVisitorImpl struct {
    db             events.DB
    paymentsClient *paymentsApi.Client
    deserializer   *EventsDeserializer
}

func NewEventsVisitorImpl(db events.DB, client *paymentsApi.Client) *EventsVisitorImpl {
    return &EventsVisitorImpl{
        db:             db,
        paymentsClient: client,
        deserializer:   NewEventsDeserializer(),
    }
}

// Acá hay que ir a payments a patchear el withdrawal
func (v *EventsVisitorImpl) VisitOutboundPaymentCreated(outboundPaymentCreated OutboundPaymentCreated) error {
    err := NewOutboundPaymentCreatedHandler(v.paymentsClient).Handle(outboundPaymentCreated)
    if err != nil {
        return fmt.Errorf("failed processing OutboundPaymentCreated: %w", err)
    }
    logger, err := zap.NewDevelopment()
    if err != nil {
        return fmt.Errorf("failed creating logger: %w", err)
    }
    logger.Info("OutboundPaymentCreated event processed successfully", zap.String("withdrawal_id", outboundPaymentCreated.WithdrawalId.String()))
    return nil
}

func (v *EventsVisitorImpl) VisitOutboundPaymentUpdated(outboundPaymentUpdated OutboundPaymentUpdated) error {
    NewOutboundPaymentUpdatedHandler(v.db, v.paymentsClient).Handle(outboundPaymentUpdated)
    return nil
}
