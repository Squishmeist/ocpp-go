package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	args := os.Args[1:]
	payload := `[2, "uuid-1", "Heartbeat", {}]`

	switch args[0] {
	case "requesterror":
		payload = `[2, "uuid-error", "Heartbeat", { "currentTimee": "2025-07-22T11:25:25.230Z" }]`
	case "confirmationerror":
		payload = `[3, "uuid-error", { "currentTimee": "2025-07-22T11:25:25.230Z" }]`
	case "heartbeatrequest":
		payload = `[2, "uuid-1", "Heartbeat", {}]`
	case "heartbeatconfirmation":
		payload = `[3, "uuid-1", { "currentTime": "2025-07-22T11:25:25.230Z" }]`
	case "bootnotificationrequest":
		payload = `[2,"uuid-2", "BootNotification",{
            "chargeBoxSerialNumber": "91234567",
            "chargePointModel": "Zappi",
            "chargePointSerialNumber": "91234567",
            "chargePointVendor": "Myenergi",
            "firmwareVersion": "5540",
            "iccid": "",
            "imsi": "",
            "meterType": "",
            "meterSerialNumber": "91234567"
        }]`
	case "bootnotificationconfirmation":
		payload = `[3,"uuid-2",{
            "currentTime": "2024-04-02T11:44:38Z",
            "interval": 30,
            "status": "Accepted"
        }]`
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
