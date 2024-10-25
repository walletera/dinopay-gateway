package app

import "log/slog"

type Option func(app *App)

func WithRabbitmqHost(host string) func(a *App) {
    return func(a *App) {
        a.rabbitmqHost = host
    }
}

func WithRabbitmqPort(port int) func(a *App) {
    return func(a *App) {
        a.rabbitmqPort = port
    }
}

func WithRabbitmqUser(user string) func(a *App) {
    return func(a *App) {
        a.rabbitmqUser = user
    }
}

func WithRabbitmqPassword(password string) func(a *App) {
    return func(a *App) {
        a.rabbitmqPassword = password
    }
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
