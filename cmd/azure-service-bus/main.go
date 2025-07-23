package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
    args := os.Args[1:]
    payload := `[2, "a9ea539a-e0b9-4d39-83e7-e40aa5b476d4", "Heartbeat", {}]`

    switch(args[0]) {
    case "error":
        payload = `[3, "a9ea539a-e0b9-4d39-83e7-e40aa5b476d4", { "currentTimee": "2025-07-22T11:25:25.230Z" }]`
    case "heartbeatrequest":
        payload = `[2, "a9ea539a-e0b9-4d39-83e7-e40aa5b476d4", "Heartbeat", {}]`
    case "heartbeatconfirmation":
        payload = `[3, "a9ea539a-e0b9-4d39-83e7-e40aa5b476d4", { "currentTime": "2025-07-22T11:25:25.230Z" }]`
    }

    ctx := context.Background()
    connStr := "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true;"
    client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
    if err != nil {
        panic(err)
    }

    topicSender, err := client.NewSender("topic.1", nil)
    if err != nil {
        panic(err)
    }
    defer topicSender.Close(ctx)

    msg := &azservicebus.Message{
        Body: []byte(payload),
    }
    if err := topicSender.SendMessage(ctx, msg, nil); err != nil {
        panic(err)
    }
    fmt.Println("Message sent to topic!")
}