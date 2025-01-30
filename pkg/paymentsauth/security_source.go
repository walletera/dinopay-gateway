package paymentsauth

import (
    "context"

    paymentsapi "github.com/walletera/payments-types/api"
)

type SecuritySource struct {
}

func NewSecuritySource() *SecuritySource {
    return &SecuritySource{}
}

func (s *SecuritySource) BearerAuth(ctx context.Context, operationName paymentsapi.OperationName) (paymentsapi.BearerAuth, error) {
    // TODO
    return paymentsapi.BearerAuth{Token: "some.dummy.token"}, nil
}
