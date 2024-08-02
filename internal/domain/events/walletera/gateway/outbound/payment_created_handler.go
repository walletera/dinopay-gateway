package outbound

import (
    "context"

    paymentsApi "github.com/walletera/payments/api"
)

type PaymentCreatedHandler struct {
    client *paymentsApi.Client
}

func NewOutboundPaymentCreatedHandler(client *paymentsApi.Client) *PaymentCreatedHandler {
    return &PaymentCreatedHandler{
        client: client,
    }
}

func (h *PaymentCreatedHandler) Handle(ctx context.Context, outboundPaymentCreated PaymentCreated) error {
    return updateWithdrawalStatus(
        ctx,
        h.client,
        outboundPaymentCreated.WithdrawalId,
        outboundPaymentCreated.DinopayPaymentId,
        outboundPaymentCreated.DinopayPaymentStatus,
    )
}
