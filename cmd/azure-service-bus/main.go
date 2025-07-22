package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
    connStr := "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true;"
    client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
    if err != nil {
        panic(err)
    }

    // Send a message to the topic
    topicSender, err := client.NewSender("topic.1", nil)
    if err != nil {
        panic(err)
    }
    defer topicSender.Close(context.Background())

    msg := &azservicebus.Message{
        Body: []byte("Hello from Go topic!"),
    }
    if err := topicSender.SendMessage(context.Background(), msg, nil); err != nil {
        panic(err)
    }
    fmt.Println("Message sent to topic!")

    // Receive the message from the subscription
    subReceiver, err := client.NewReceiverForSubscription("topic.1", "subscription.1", nil)
    if err != nil {
        panic(err)
    }
    defer subReceiver.Close(context.Background())

    messages, err := subReceiver.ReceiveMessages(context.Background(), 1, nil)
    if err != nil {
        panic(err)
    }
    for _, msg := range messages {
        fmt.Printf("Received message from subscription: %s\n", string(msg.Body))
        subReceiver.CompleteMessage(context.Background(), msg, nil)
    }
}