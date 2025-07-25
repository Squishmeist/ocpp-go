# ⚡️ ocpp-go

A Go implementation for handling OCPP 1.6 messages, designed for extensible server-side processing and integration with Azure Service Bus.

## 🛠️ How it Works

- 📨 **OCPP Listener:** Subscribes to a local Azure Service Bus topic and processes incoming OCPP messages.
- 🟦 **Azure Service Bus:** Provides a local, Dockerized message bus for development and testing.
- 📤 **HTTP Sender Server:** Lets you send OCPP messages to the topic via HTTP for easy testing.
- 👀 **Topic Receivers:** Automatically logs all messages received on both inbound and outbound topics.

## ✨ Features

- ✅ Parse and validate OCPP 1.6 messages (Heartbeat, BootNotification, etc.)
- ☁️ Azure Service Bus integration (queue and topic support)
- 🗃️ In-memory state management and message pairing
- 🔎 Real-time logging of all topic traffic

## 🚀 Getting Started

### 🛠️ Development Workflow

1. **Start the Azure Service Bus emulator (Docker Compose):**

   ```sh
   make azure-service-bus
   ```

   Runs the emulator in Docker using your `docker-compose.yaml`.

2. **Start the OCPP listener:**

   ```sh
   make ocpp
   ```

   Runs the OCPP machine, which listens for messages on the inbound topic and processes them.

3. **Start the HTTP sender server:**
   ```sh
   make send-message
   ```
   Starts an HTTP server for sending OCPP messages to the inbound topic.  
   Also runs receivers for both inbound and outbound topics, so you can see all message traffic in your logs.

### 🟦 Azure Service Bus

A local Azure Service Bus emulator for development.

> **Note:** Make sure the emulator is running before starting the OCPP listener or sender server.

### ⚡️ OCPP

The OCPP machine listens for messages from your local Azure Service Bus inbound topic, parses, and processes them.

### 📤 Send-Messages

Send OCPP messages to your local Azure Service Bus topic using the HTTP server:

```sh
curl -X POST "http://localhost:<YOUR-HTTP-PORT>/send?msg=heartbeatrequest"
```

**Supported `msg` values:**

- `heartbeatrequest`
- `heartbeatconfirmation`
- `bootnotificationrequest`
- `bootnotificationconfirmation`
- `requesterror`
- `confirmationerror`

Each value triggers a different OCPP payload.

### 👀 Message Receivers

The sender server runs receivers for both the inbound and outbound topics.  
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
