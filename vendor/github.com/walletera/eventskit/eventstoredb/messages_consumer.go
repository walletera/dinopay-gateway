package eventstoredb

import (
    "context"
    "fmt"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/walletera/eventskit/messages"
)

type MessagesConsumer struct {
    esdbClient *esdb.Client

    connectionString string
    streamName       string
    groupName        string
}

func NewMessagesConsumer(connectionString, streamName string, groupName string, opts ...Opt) (*MessagesConsumer, error) {
    messagesConsumer := &MessagesConsumer{
        streamName: streamName,
        groupName:  groupName,
    }

    for _, opt := range opts {
        opt(messagesConsumer)
    }

    esdbClient, err := GetESDBClient(connectionString)
    if err != nil {
        return nil, err
    }

    messagesConsumer.esdbClient = esdbClient
    return messagesConsumer, nil
}

func (mc *MessagesConsumer) Consume() (<-chan messages.Message, error) {
    persistentSubscription, err := mc.esdbClient.SubscribeToPersistentSubscription(
        context.Background(),
        mc.streamName,
        mc.groupName,
        esdb.SubscribeToPersistentSubscriptionOptions{},
    )
    if err != nil {
        panic(err)
    }
    messagesCh := make(chan messages.Message)
    go func() {
        defer close(messagesCh)
        for {
            persistentSubscriptionEvent := persistentSubscription.Recv()
            if persistentSubscriptionEvent.SubscriptionDropped != nil {
                fmt.Printf("persistent subscription dropped: %s", persistentSubscriptionEvent.SubscriptionDropped.Error.Error())
                return
            }
            event := persistentSubscriptionEvent.EventAppeared.Event
            originalEvent := event.Event
            if originalEvent == nil {
                fmt.Printf("original eventAppeared is nil in persistent subcription eventAppeared")
                return
            }
            messagesCh <- messages.NewMessage(
                originalEvent.Data,
                NewAcknowledger(persistentSubscription, persistentSubscriptionEvent.EventAppeared),
            )
        }
    }()
    return messagesCh, nil
}

func (mc *MessagesConsumer) Close() error {
    err := mc.esdbClient.Close()
    if err != nil {
        return fmt.Errorf("failed closing eventstoredb message consumer: %w", err)
    }
    return nil
}
