# ⚡️ ocpp-go

A Go implementation for handling OCPP 1.6 messages, designed for extensible server-side processing and integration with Azure Service Bus.

## 🛠️ How it Works

- 📨 The OCPP listener subscribes to a local Azure Service Bus topic and processes incoming OCPP messages.
- 🟦 The Azure Service Bus Emulator provides a local, Dockerized message bus for development.
- 🛠️ Use the CLI to send OCPP messages to the topic for testing and development.

## ✨ Features

- ✅ Parse and validate OCPP 1.6 messages (Heartbeat, BootNotification, etc.)
- ☁️ Azure Service Bus integration (queue and topic support)
- 🗃️ In-memory state management and message pairing

## 🚀 Getting Started

### 🛠️ Development Workflow

1. **Start the Azure Service Bus emulator:**
   ```sh
   make azure-service-bus
   ```
2. **Start the OCPP listener:**
   ```sh
   make ocpp
   ```
3. **Send a message to the topic:**
   ```sh
   make send-message ARGS=heartbeatrequest
   ```

### ⚡️ OCPP Listener

To start the OCPP listener service (which listens for messages from your local Azure Service Bus topic):

```sh
make ocpp
```

> **Note:** Make sure the Azure Service Bus emulator is running before starting the listener.

### 🟦 Azure Service Bus Emulator

To run the Azure Service Bus emulator locally (for topic/queue development):

```sh
make azure-service-bus
```

### 📤 Sending Messages

You can send OCPP messages to your local Azure Service Bus topic using the CLI:

```sh
make send-message ARGS=heartbeatrequest
```

Supported ARGS:

- `heartbeatrequest`
- `heartbeatconfirmation`
- (add more as needed)

### 📦 Payloads

OCPP messages are handled as arrays, mapped to these Go structs:

```go
type RequestBody struct {
    MessageType MessageType // 2
    MessageId   string      // UUID
    Action      ActionType  // e.g. Heartbeat
    Payload     []byte      // JSON-encoded payload
}

type ConfirmationBody struct {
    MessageType MessageType // 3
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
