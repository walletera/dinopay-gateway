package events

type Handler interface {
    HandleEvent(event any) error
}
