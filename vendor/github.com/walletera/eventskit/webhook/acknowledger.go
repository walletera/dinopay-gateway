package webhook

import (
    "net/http"

    "github.com/walletera/eventskit/messages"
    "github.com/walletera/werrors"
)

type Acknowledger struct {
    done   chan any
    writer http.ResponseWriter
}

func NewAcknowledger(writer http.ResponseWriter) *Acknowledger {
    return &Acknowledger{
        done:   make(chan any),
        writer: writer,
    }
}

func (a *Acknowledger) Ack() error {
    a.writer.WriteHeader(http.StatusCreated)
    close(a.done)
    return nil
}

func (a *Acknowledger) Nack(opts messages.NackOpts) error {
    var status int
    switch opts.ErrorCode {
    case werrors.UnprocessableMessageErrorCode:
        status = http.StatusBadRequest
    default:
        status = http.StatusInternalServerError
    }
    a.writer.WriteHeader(status)
    close(a.done)
    return nil
}

func (a *Acknowledger) Done() <-chan any {
    return a.done
}
