package app

import "log/slog"

type Option func(app *App)

func WithRabbitMQUrl(url string) func(app *App) {
    return func(app *App) { app.rabbitmqUrl = url }
}

func WithDinopayUrl(url string) func(app *App) {
    return func(app *App) { app.dinopayUrl = url }
}

func WithPaymentsUrl(url string) func(app *App) {
    return func(app *App) { app.paymentsUrl = url }
}

func WithESDBUrl(url string) func(app *App) {
    return func(app *App) { app.esdbUrl = url }
}

func WithLogHandler(handler slog.Handler) func(app *App) {
    return func(app *App) { app.logHandler = handler }
}
