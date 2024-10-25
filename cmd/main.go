package main

import (
    "context"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/walletera/dinopay-gateway/internal/app"
)

const shutdownTimeout = 10 * time.Second

func main() {
    ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer ctxCancel()

    rabbitmqHost := mustGetEnv("RABBITMQ_HOST")
    rabbitmqPort := mustGetIntEnv("RABBITMQ_PORT")
    rabbitmqUser := mustGetEnv("RABBITMQ_USER")
    rabbitmqPassword := mustGetEnv("RABBITMQ_PASSWORD")
    dinopayURL := mustGetEnv("DINOPAY_URL")
    paymentsURL := mustGetEnv("PAYMENTS_URL")
    eventstoredbURL := mustGetEnv("EVENTSTOREDB_URL")

    app, err := app.NewApp(
        app.WithRabbitmqHost(rabbitmqHost),
        app.WithRabbitmqPort(rabbitmqPort),
        app.WithRabbitmqUser(rabbitmqUser),
        app.WithRabbitmqPassword(rabbitmqPassword),
        app.WithDinopayUrl(dinopayURL),
        app.WithPaymentsUrl(paymentsURL),
        app.WithESDBUrl(eventstoredbURL),
    )
    if err != nil {
        panic(err)
    }

    err = app.Run(ctx)
    if err != nil {
        panic(err)
    }

    <-ctx.Done()

    shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), shutdownTimeout)
    defer shutdownCtxCancel()

    app.Stop(shutdownCtx)
}

func mustGetEnv(envName string) string {
    value, found := os.LookupEnv(envName)
    if !found {
        panic("env var not defined: " + envName)
    }
    return value
}

func mustGetIntEnv(envName string) int {
    strEnvValue := mustGetEnv(envName)
    intEnvValue, err := strconv.Atoi(strEnvValue)
    if err != nil {
        panic("env var is not an int: " + envName)
    }
    return intEnvValue
}
