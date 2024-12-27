package webhook

import (
    "io"
    "log/slog"
    "net/http"
    "time"

    "github.com/walletera/eventskit/messages"
)

const (
    MessageProcessingTimeout = 10 * time.Second
)

type handler struct {
    logger *slog.Logger
    msgCh  chan messages.Message
}

func newHandler(msgCh chan messages.Message, logger *slog.Logger) *handler {
    return &handler{logger: logger, msgCh: msgCh}
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodPost {
        writer.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    rawBody, err := io.ReadAll(request.Body)
    if err != nil {
        h.logError("failed reading request body", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }
    if rawBody == nil || len(rawBody) == 0 {
        h.logError("webhook request contains empty body", nil)
        writer.WriteHeader(http.StatusBadRequest)
        return
    }
    timeoutC := time.After(MessageProcessingTimeout + (1 * time.Second))
    acknowledger := NewAcknowledger(writer)
    h.msgCh <- messages.NewMessage(rawBody, acknowledger)
    select {
    case <-acknowledger.Done():
        return
    case <-timeoutC:
        // This is a safeguard measure, it should not happen. The processor should
        // time out and cancel the operation before us.
        h.logError("webhook server timeout waiting for message to be processed (should not had happened)", nil)
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func (h *handler) logError(msg string, err error) {
    if h.logger == nil {
        return
    }
    if err != nil {
        h.logger.Error(msg, slog.String("error", err.Error()))
    } else {
        h.logger.Error(msg)
    }
}

func (h *handler) logDebug(msg string) {
    if h.logger == nil {
        return
    }
    h.logger.Debug(msg)
}
