package dinopay

import (
    "context"

    "github.com/walletera/dinopay/api"
)

type Client interface {
    CreatePayment(ctx context.Context, req *api.Payment) (api.CreatePaymentRes, error)
}
