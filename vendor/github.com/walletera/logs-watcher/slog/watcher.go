package slog

import (
    "log/slog"

    "github.com/walletera/logs-watcher/newline"
)

// Watcher watch logs produced by
// a slog logger
type Watcher struct {
    *newline.Watcher
    decoratedHandler slog.Handler
}

func NewWatcher(handler slog.Handler) *Watcher {
    watcher := newline.NewWatcher()
    logHook := func(record slog.Record) {
        watcher.AddLogLine(record.Message)
    }
    decoratedHandler := NewHandlerDecorator(handler, logHook)
    return &Watcher{
        Watcher:          watcher,
        decoratedHandler: decoratedHandler,
    }
}

func (s *Watcher) DecoratedHandler() slog.Handler {
    return s.decoratedHandler
}
