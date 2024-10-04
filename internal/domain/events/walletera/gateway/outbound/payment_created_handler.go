package outbound

import (
    "context"

    paymentsapi "github.com/walletera/payments-types/api"
)

type PaymentCreatedHandler struct {
    client *paymentsapi.Client
}

func NewOutboundPaymentCreatedHandler(client *paymentsapi.Client) *PaymentCreatedHandler {
    return &PaymentCreatedHandler{
        client: client,
    }
}

func (h *PaymentCreatedHandler) Handle(ctx context.Context, outboundPaymentCreated PaymentCreated) error {
    return updatePaymentStatus(
        ctx,
        h.client,
        outboundPaymentCreated.WithdrawalId,
        outboundPaymentCreated.DinopayPaymentId,
        outboundPaymentCreated.DinopayPaymentStatus,
    )
}
