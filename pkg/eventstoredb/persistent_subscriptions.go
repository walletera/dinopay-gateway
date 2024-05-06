package eventstoredb

import (
    "context"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
)

func CreatePersistentSubscription(connectionString string, streamName string, groupName string) error {
    settings, err := esdb.ParseConnectionString(connectionString)
    if err != nil {
        return err
    }

    esdbClient, err := esdb.NewClient(settings)
    if err != nil {
        return err
    }

    subscriptionSettings := esdb.SubscriptionSettingsDefault()
    subscriptionSettings.ResolveLinkTos = true
    subscriptionSettings.MaxRetryCount = 3

    err = esdbClient.CreatePersistentSubscription(
        context.Background(),
        streamName,
        groupName,
        esdb.PersistentStreamSubscriptionOptions{
            Settings: &subscriptionSettings,
        },
    )
    if err != nil {
        println("CreatePersistentSubscription failed: " + err.Error())
    }

    return nil
}
