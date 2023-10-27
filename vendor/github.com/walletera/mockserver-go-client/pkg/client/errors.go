package client

import "fmt"

type IncorrectRequestFormat struct{}

func (i *IncorrectRequestFormat) Error() string {
    return "incorrect request format"
}

type InvalidExpectation struct{}

func (i *InvalidExpectation) Error() string {
    return "invalid expectation"
}

type RequestHasNotBeenReceived struct{}

func (r *RequestHasNotBeenReceived) Error() string {
    return "request has not been received specified numbers of times"
}

type UnexpectedStatusCode struct {
    Endpoint   string
    StatusCode int
}

func (u *UnexpectedStatusCode) Error() string {
    return fmt.Sprintf("unexpected status code %d received from endpoint %s", u.StatusCode, u.Endpoint)
}
