package gateway

import (
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

func (h *OutboundPaymentCreatedHandler) Handle(outboundPaymentCreated OutboundPaymentCreated) error {
    return updateWithdrawalStatus(
        h.client,
        outboundPaymentCreated.WithdrawalId,
        outboundPaymentCreated.DinopayPaymentId,
        outboundPaymentCreated.DinopayPaymentStatus,
    )
}
