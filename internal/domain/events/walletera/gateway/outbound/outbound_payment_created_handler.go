package outbound

import (
    "context"

    paymentsApi "github.com/walletera/payments/api"
)

type OutboundPaymentCreatedHandler struct {
    client *paymentsApi.Client
}

func NewOutboundPaymentCreatedHandler(client *paymentsApi.Client) *OutboundPaymentCreatedHandler {
    return &OutboundPaymentCreatedHandler{
        client: client,
    }
}

func (h *OutboundPaymentCreatedHandler) Handle(ctx context.Context, outboundPaymentCreated OutboundPaymentCreated) error {
    return updateWithdrawalStatus(
        ctx,
        h.client,
        outboundPaymentCreated.WithdrawalId,
        outboundPaymentCreated.DinopayPaymentId,
        outboundPaymentCreated.DinopayPaymentStatus,
    )
}
