// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateEventSubscription implements createEventSubscription operation.
	//
	// Subscribe to events.
	//
	// POST /webhooks/subscriptions
	CreateEventSubscription(ctx context.Context, req *EventSubscription) error
	// CreatePayment implements createPayment operation.
	//
	// Create a new payment.
	//
	// POST /payments
	CreatePayment(ctx context.Context, req *Payment) (CreatePaymentRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
