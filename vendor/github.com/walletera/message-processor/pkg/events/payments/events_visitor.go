package payments

type EventsVisitor interface {
    VisitWithdrawalCreated(withdrawalCreated WithdrawalCreated) error
}
