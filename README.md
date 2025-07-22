# ⚡️ ocpp-go

A Go implementation for handling OCPP 1.6 messages, designed for extensible server-side processing and integration with Azure Service Bus.

## 📁 Structure

- `service/ocpp/` — OCPP message types, utilities, handlers
- `service/azure-service-bus` — Docker Azure Service Bus Emulator

## ✨ Features

- ✅ Parse and validate OCPP 1.6 messages (Heartbeat, BootNotification, etc.)
- ☁️ Azure Service Bus integration (queue and topic support)
- 🗃️ In-memory state management and message pairing

## 🚀 Getting Started

### OCPP Development

For development, an HTTP server is run and POST requests can be sent to `/test`:

```sh
make ocpp
```
or
```sh
go run ./cmd/ocpp/main.go
```

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
[
    2,
    "uuid-1",
    "Heartbeat",
    {}
]
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