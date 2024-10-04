package events

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type Handler interface {
    HandlePaymentCreated(ctx context.Context, paymentCreatedEvent PaymentCreated) errors.ProcessingError
    HandlePaymentUpdated(ctx context.Context, paymentCreatedEvent PaymentUpdated) errors.ProcessingError
}
