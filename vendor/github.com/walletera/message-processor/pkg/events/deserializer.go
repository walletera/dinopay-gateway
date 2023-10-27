package events

type Deserializer[Visitor any] interface {
    Deserialize(rawEvent []byte) (Event[Visitor], error)
}
