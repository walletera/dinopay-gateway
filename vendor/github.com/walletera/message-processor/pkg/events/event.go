package events

type Event[Visitor any] interface {
    Accept(visitor Visitor) error
}
