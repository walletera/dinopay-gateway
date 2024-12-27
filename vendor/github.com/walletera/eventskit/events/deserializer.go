package events

type Deserializer[Handler any] interface {
    Deserialize(rawEvent []byte) (Event[Handler], error)
}
