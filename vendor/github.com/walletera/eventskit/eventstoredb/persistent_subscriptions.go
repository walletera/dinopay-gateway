package eventstoredb

import (
    "context"
    "errors"
    "fmt"

    "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

func CreatePersistentSubscription(connectionString string, streamName string, groupName string) error {
    esdbClient, err := GetESDBClient(connectionString)
    if err != nil {
        return err
    }

    subscriptionSettings := kurrentdb.SubscriptionSettingsDefault()
    subscriptionSettings.ResolveLinkTos = true

    err = esdbClient.CreatePersistentSubscription(
        context.Background(),
        streamName,
        groupName,
        kurrentdb.PersistentStreamSubscriptionOptions{
            Settings: &subscriptionSettings,
        },
    )
    if err != nil {
        var esdbError *kurrentdb.Error
        ok := errors.As(err, &esdbError)
        if !ok || !esdbError.IsErrorCode(kurrentdb.ErrorCodeResourceAlreadyExists) {
            return fmt.Errorf("failed creating persistent subscription for stream %s and group %s: %w", streamName, groupName, err)
        }
    }

    return nil
}
