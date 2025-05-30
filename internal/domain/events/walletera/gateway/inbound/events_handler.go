package inbound

import (
    "context"
    "log/slog"

    "github.com/walletera/dinopay-gateway/pkg/logattr"
    builders "github.com/walletera/payments-types/builders/privateapi"
    paymentsapi "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

type EventsHandler interface {
    HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError
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

func (ev *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError {
    postPaymentReq := &paymentsapi.PostPaymentReq{
        ID:         inboundPaymentReceived.PaymentId,
        Amount:     inboundPaymentReceived.Amount,
        Currency:   paymentsapi.Currency(inboundPaymentReceived.Currency),
        Gateway:    paymentsapi.GatewayDinopay,
        Direction:  paymentsapi.DirectionInbound,
        CustomerId: inboundPaymentReceived.CustomerId,
        Status:     paymentsapi.PaymentStatusConfirmed,
        ExternalId: paymentsapi.NewOptString(inboundPaymentReceived.DinopayPaymentId.String()),
        Debtor: builders.NewDinopayAccountBuilder().
            WithCurrency(paymentsapi.Currency(inboundPaymentReceived.Currency)).
            WithAccountHolder(inboundPaymentReceived.SourceAccount.AccountHolder).
            WithAccountNumber(inboundPaymentReceived.SourceAccount.AccountNumber).
            Build(),
        Beneficiary: builders.NewDinopayAccountBuilder().
            WithCurrency(paymentsapi.Currency(inboundPaymentReceived.Currency)).
            WithAccountHolder(inboundPaymentReceived.DestinationAccount.AccountHolder).
            WithAccountNumber(inboundPaymentReceived.DestinationAccount.AccountNumber).
            Build(),
    }
    _, err := ev.paymentsApiClient.PostPayment(ctx, postPaymentReq, paymentsapi.PostPaymentParams{})
    if err != nil {
        // TODO handle this error properly
        ev.logger.Error("failed creating payment on payments api", logattr.Error(err.Error()))
        return werrors.NewRetryableInternalError(err.Error())
    }
    // TODO handle response
    ev.logger.Info("Gateway event InboundPaymentReceived processed successfully", logattr.EventType(inboundPaymentReceived.Type()))
    return nil
}
