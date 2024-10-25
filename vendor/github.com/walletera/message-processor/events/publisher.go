package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, data EventData, topic string) error
}
