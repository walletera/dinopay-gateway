package outbound

import (
    "context"

    paymentsapi "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

type PaymentCreatedHandler struct {
    client *paymentsapi.Client
}

func NewOutboundPaymentCreatedHandler(client *paymentsapi.Client) *PaymentCreatedHandler {
    return &PaymentCreatedHandler{
        client: client,
    }
}

func (h *PaymentCreatedHandler) Handle(ctx context.Context, outboundPaymentCreated PaymentCreated) werrors.WError {
    return updatePaymentStatus(
        ctx,
        h.client,
        outboundPaymentCreated.PaymentId,
        outboundPaymentCreated.DinopayPaymentId,
        outboundPaymentCreated.DinopayPaymentStatus,
    )
}
