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

    webhookPort := mustGetIntEnv("WEBHOOK_PORT")
    rabbitmqHost := mustGetEnv("RABBITMQ_HOST")
    rabbitmqPort := mustGetIntEnv("RABBITMQ_PORT")
    rabbitmqUser := mustGetEnv("RABBITMQ_USER")
    rabbitmqPassword := mustGetEnv("RABBITMQ_PASSWORD")
    dinopayURL := mustGetEnv("DINOPAY_URL")
    accountsURL := mustGetEnv("ACCOUNTS_URL")
    paymentsURL := mustGetEnv("PAYMENTS_URL")
    eventstoredbURL := mustGetEnv("EVENTSTOREDB_URL")

    dinopayGateway, err := app.NewApp(
        app.WithWebhookPort(webhookPort),
        app.WithRabbitmqHost(rabbitmqHost),
        app.WithRabbitmqPort(rabbitmqPort),
        app.WithRabbitmqUser(rabbitmqUser),
        app.WithRabbitmqPassword(rabbitmqPassword),
        app.WithDinopayUrl(dinopayURL),
        app.WithAccountsUrl(accountsURL),
        app.WithPaymentsUrl(paymentsURL),
        app.WithESDBUrl(eventstoredbURL),
    )
    if err != nil {
        panic(err)
    }

    err = dinopayGateway.Run(ctx)
    if err != nil {
        panic(err)
    }

    <-ctx.Done()

    shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), shutdownTimeout)
    defer shutdownCtxCancel()

    dinopayGateway.Stop(shutdownCtx)
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
