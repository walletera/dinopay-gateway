package eventstoredb

import (
    "os"

    "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

func GetESDBClient(connectionString string) (*kurrentdb.Client, error) {
    settings, err := kurrentdb.ParseConnectionString(connectionString)
    if err != nil {
        return nil, err
    }

    value, ok := os.LookupEnv("EVENTSTOREDB_CLIENT_LOG_ENABLED")
    if !ok || value != "true" {
        settings.Logger = func(level kurrentdb.LogLevel, format string, args ...interface{}) {}
    }

    return kurrentdb.NewClient(settings)
}
