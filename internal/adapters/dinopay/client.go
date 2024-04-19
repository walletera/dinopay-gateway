package dinopay

import (
    "context"
    "fmt"
    "net/http"

    "github.com/walletera/dinopay/api"
)

type Client struct {
    client *api.Client
}

func NewClient(url string) (*Client, error) {
    client, err := api.NewClient(url, api.WithClient(http.DefaultClient))
    if err != nil {
        return nil, fmt.Errorf("failed creating dinopay api client: %w", err)
    }
    return &Client{
        client: client,
    }, nil
}

func (c *Client) CreatePayment(ctx context.Context, req *api.Payment) (api.CreatePaymentRes, error) {
    return c.client.CreatePayment(ctx, req)
}
