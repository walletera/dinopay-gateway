package client

import "net/http"

type IHTTPClient interface {
    Do(r *http.Request) (*http.Response, error)
}
