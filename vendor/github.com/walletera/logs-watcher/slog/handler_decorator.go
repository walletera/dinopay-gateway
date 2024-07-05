package slog

import (
    "context"
    "log/slog"
)

type HandlerDecorator struct {
    handler slog.Handler
    logHook LogHook
}

type LogHook func(record slog.Record)

func NewHandlerDecorator(handler slog.Handler, hook LogHook) *HandlerDecorator {
    return &HandlerDecorator{
        handler: handler,
        logHook: hook,
    }
}

func (h *HandlerDecorator) Enabled(ctx context.Context, level slog.Level) bool {
    return h.handler.Enabled(ctx, level)
}

func (h *HandlerDecorator) Handle(ctx context.Context, record slog.Record) error {
    h.logHook(record)
    return h.handler.Handle(ctx, record)
}

func (h *HandlerDecorator) WithAttrs(attrs []slog.Attr) slog.Handler {
    h.handler = h.handler.WithAttrs(attrs)
    return h
}

func (h *HandlerDecorator) WithGroup(name string) slog.Handler {
    h.handler = h.handler.WithGroup(name)
    return h
}
