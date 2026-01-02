package outbound

import (
	"context"
	"log/slog"

	"github.com/walletera/dinopay-gateway/pkg/logattr"
	paymentsapi "github.com/walletera/payments-types/privateapi"
	"github.com/walletera/werrors"
)

type PaymentCreatedHandler struct {
	client *paymentsapi.Client
	logger *slog.Logger
}

func NewOutboundPaymentCreatedHandler(client *paymentsapi.Client, logger *slog.Logger) *PaymentCreatedHandler {
	return &PaymentCreatedHandler{
		client: client,
		logger: logger,
	}
}

func (h *PaymentCreatedHandler) Handle(ctx context.Context, outboundPaymentCreated PaymentCreated) werrors.WError {
	h.logger.Info("handling OutboundPaymentCreated event", logattr.PaymentId(outboundPaymentCreated.PaymentId.String()))
	return updatePaymentStatus(
		ctx,
		h.client,
		outboundPaymentCreated.PaymentId,
		outboundPaymentCreated.DinopayPaymentId,
		outboundPaymentCreated.DinopayPaymentStatus,
	)
}
