package messages

type Consumer interface {
    Consume() (<-chan Message, error)
    Close() error
}
