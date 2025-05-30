package eventstoredb

import (
    "context"
    "errors"
    "fmt"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
)

func CreatePersistentSubscription(connectionString string, streamName string, groupName string, subscriptionSettings esdb.PersistentSubscriptionSettings) error {
    esdbClient, err := GetESDBClient(connectionString)
    if err != nil {
        return err
    }

    err = esdbClient.CreatePersistentSubscription(
        context.Background(),
        streamName,
        groupName,
        esdb.PersistentStreamSubscriptionOptions{
            Settings: &subscriptionSettings,
        },
    )
    if err != nil {
        var esdbError *esdb.Error
        ok := errors.As(err, &esdbError)
        if !ok || !esdbError.IsErrorCode(esdb.ErrorCodeResourceAlreadyExists) {
            return fmt.Errorf("failed creating persistent subscription for stream %s and group %s: %w", streamName, groupName, err)
        }
    }

    return nil
}
