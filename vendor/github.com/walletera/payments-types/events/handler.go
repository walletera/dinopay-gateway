package events

import (
    "context"

    "github.com/walletera/werrors"
)

type Handler interface {
    HandlePaymentCreated(ctx context.Context, paymentCreatedEvent PaymentCreated) werrors.WError
    HandlePaymentUpdated(ctx context.Context, paymentCreatedEvent PaymentUpdated) werrors.WError
}
