package eventstoredb

import (
    "context"
    "fmt"

    "github.com/EventStore/EventStore-Client-Go/v4/esdb"
    "github.com/walletera/message-processor/messages"
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
        for {
            persistentSubscriptionEvent := persistentSubscription.Recv()
            if persistentSubscriptionEvent.SubscriptionDropped != nil {
                fmt.Printf("persistent subscription dropped: %s", persistentSubscriptionEvent.SubscriptionDropped.Error.Error())
                return
            }
            event := persistentSubscriptionEvent.EventAppeared.Event
            originalEvent := event.Event
            if originalEvent == nil {
                fmt.Printf("original event is nil in persistent subcription event")
                return
            }
            // TODO The Ack/Nack must be done in the MesssageProcessor
            persistentSubscription.Ack(event)
            messagesCh <- messages.Message{
                Payload: originalEvent.Data,
            }
        }
    }()
    return messagesCh, nil
}

func (mc *MessagesConsumer) Close() error {
    //TODO implement me
    panic("implement me")
}