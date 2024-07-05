package payments

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type EventsVisitor interface {
    VisitWithdrawalCreated(ctx context.Context, withdrawalCreated WithdrawalCreatedEvent) errors.ProcessingError
}
