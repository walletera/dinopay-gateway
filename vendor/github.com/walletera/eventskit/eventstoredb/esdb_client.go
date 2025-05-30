package eventstoredb

import (
    "os"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
)

func GetESDBClient(connectionString string) (*esdb.Client, error) {
    settings, err := esdb.ParseConnectionString(connectionString)
    if err != nil {
        return nil, err
    }

    value, ok := os.LookupEnv("EVENTSTOREDB_CLIENT_LOG_ENABLED")
    if !ok || value != "true" {
        settings.Logger = func(level esdb.LogLevel, format string, args ...interface{}) {}
    }

    return esdb.NewClient(settings)
}
