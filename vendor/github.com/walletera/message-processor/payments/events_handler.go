package payments

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type EventsHandler interface {
    HandleWithdrawalCreated(ctx context.Context, withdrawalCreated WithdrawalCreatedEvent) errors.ProcessingError
}
