package events

import "context"

type RoutingInfo struct {
    Topic      string
    RoutingKey string
}

type Publisher interface {
    Publish(ctx context.Context, data EventData, info RoutingInfo) error
}
