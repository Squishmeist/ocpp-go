# ⚡️ ocpp-go

A Go implementation for handling OCPP 1.6 messages, designed for extensible server-side processing and integration with Azure Service Bus.

## 🛠️ How it Works

- 📨 **OCPP:** Subscribes to a local Azure Service Bus topic and processes incoming OCPP messages.
- 🟦 **Azure Service Bus:** Provides a local, Dockerized message bus for development and testing.
- 📤 **Message:** Lets you send OCPP messages to the topic via gRPC for easy testing.
- 👀 **Topic Receivers:** Automatically logs all messages received on both inbound and outbound topics.

## ✨ Features

- ✅ Parse and validate OCPP 1.6 messages (Heartbeat, BootNotification, etc.)
- ☁️ Azure Service Bus integration (queue and topic support)
- 🗃️ State management and message pairing
- 🔎 Real-time logging of all topic traffic

## 🚀 Getting Started

### 🛠️ Development Workflow

1. **Start the Azure Service Bus emulator (Docker Compose):**

   ```sh
   make azure-service-bus
   ```

   Runs the emulator in Docker using your `docker-compose.yaml`.

2. **Start the OCPP machine:**

   ```sh
   make ocpp
   ```

   Runs the OCPP machine, which listens for messages on the inbound topic and processes them.

3. **Start the Message server:**
   ```sh
   make message
   ```
   Starts a gRPC server for sending OCPP messages to the inbound topic.  
   Also runs receivers for both inbound and outbound topics, so you can see all message traffic in your logs.

### 🟦 Azure Service Bus

A local Azure Service Bus emulator for development.

> **Note:** Make sure the emulator is running before starting the OCPP machine or message server.

### ⚡️ OCPP

The OCPP machine listens for messages from your local Azure Service Bus inbound topic, parses, and processes them.

### 📤 Message

Send OCPP messages to your local Azure Service Bus topic using the gRPC server:

**Using Postman:**

1. Create a new gRPC request
2. Set server URL: `localhost:8082`
3. Select service method: `OCPPService/<message-type>`

> **💡 Tip:** The gRPC server uses reflection, so Postman will automatically discover available service methods.

**Available Message Types:**

- `HeartbeatRequest`
- `HeartbeatConfirmation`
- `BootNotificationRequest`
- `BootNotificationConfirmation`

Each triggers a different OCPP payload.

**Response Format:**

```json
{
  "message": "Message sent successfully"
}
```

### 👀 Topic Receivers

The message server runs receivers for both the inbound and outbound topics.  
This lets you observe all messages sent to either topic (from the OCPP machine or from your sender) directly in your logs.

### 📦 Payloads

OCPP messages are handled as arrays, mapped to these Go structs:

```go
type RequestBody struct {
    MessageId   string      // UUID
    Action      ActionKind  // e.g. Heartbeat
    Payload     []byte      // JSON-encoded payload
}

type ConfirmationBody struct {
    MessageId   string      // UUID
    Payload     []byte      // JSON-encoded payload
}
```

#### 📨 Examples

**RequestBody**

```json
[2, "uuid-1", "Heartbeat", {}]
```

**ConfirmationBody**

```json
[
  3,
  "uuid-1",
  {
    "currentTime": "2025-07-22T11:25:25.230Z"
  }
]
```
