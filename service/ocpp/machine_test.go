package ocpp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

type mockStore struct{}

func (m *mockStore) GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, core.HandlerResponse) {
	if uuid != "uuid-000" {
		return types.RequestBody{}, core.HandlerResponse{
			TraceID: "trace-id-123",
			Error:   utils.StringPtr("request not found"),
			Message: "request not found",
		}
	}

	return types.RequestBody{
			Uuid:    "uuid-000",
			Action:  types.BootNotification,
			Payload: []byte(`{"chargeBoxSerialNumber": "91234567"}`),
		}, core.HandlerResponse{
			TraceID: "trace-id-123",
			Error:   nil,
			Message: "request found",
		}
}

func (m *mockStore) AddRequestMessage(ctx context.Context, request types.RequestBody) (types.MessageBody, core.HandlerResponse) {
	return types.MessageBody{
			Uuid:    request.Uuid,
			Kind:    "REQUEST",
			Payload: request.Payload,
		}, core.HandlerResponse{
			TraceID: "trace-id-123",
			Message: "request added",
			Error:   nil,
		}
}

func (m *mockStore) AddConfirmationMessage(ctx context.Context, confirmation types.ConfirmationBody) (types.MessageBody, core.HandlerResponse) {
	return types.MessageBody{
			Uuid:    confirmation.Uuid,
			Kind:    "CONFIRMATION",
			Payload: confirmation.Payload,
		}, core.HandlerResponse{
			TraceID: "trace-id-456",
			Message: "confirmation added",
			Error:   nil,
		}
}

func TestHandleMessage(t *testing.T) {

	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("InvalidKind", func(t *testing.T) {
		body := []any{"2.0", "uuid-123", types.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("NoKind", func(t *testing.T) {
		body := []any{"uuid-123", types.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoAction", func(t *testing.T) {
		body := []any{2.0, "uuid-123", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_UnknownAction", func(t *testing.T) {
		body := []any{2.0, "uuid-123", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-123", types.Heartbeat}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_MissingPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-123", types.BootNotification, map[string]any{
			"chargeBoxSerialNumber":   "91234567",
			"chargePointModel":        "Zappi",
			"chargePointSerialNumber": "91234567",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("Request", func(t *testing.T) {
		body := []any{2.0, "uuid-123", types.BootNotification, map[string]any{
			"chargeBoxSerialNumber":   "91234567",
			"chargePointModel":        "Zappi",
			"chargePointSerialNumber": "91234567",
			"chargePointVendor":       "Myenergi",
			"firmwareVersion":         "5540",
			"iccid":                   "",
			"imsi":                    "",
			"meterType":               "",
			"meterSerialNumber":       "91234567",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.NoError(t, err)
	})

	t.Run("InvalidConfirmation_NoPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-123"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidConfirmation_MissingPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-000", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidConfirmation_NoMatchingRequest", func(t *testing.T) {
		body := []any{3.0, "uuid-unknown", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
			"status":      "Accepted",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.Error(t, err)
	})

	t.Run("Confirmation", func(t *testing.T) {
		body := []any{3.0, "uuid-000", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
			"status":      "Accepted",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, raw)
		assert.NoError(t, err)
	})

}

func TestParseRawMessage(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)

	t.Run("NotJSON", func(t *testing.T) {
		raw := []byte(`not a json array`)
		_, err := machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidLength", func(t *testing.T) {
		body := []any{2.0}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidKind", func(t *testing.T) {
		body := []any{"not-a-number", "uuid-789"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("UnknownKind", func(t *testing.T) {
		body := []any{99.0, "uuid-999"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoAction", func(t *testing.T) {
		body := []any{2.0, "uuid-123", 2, map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_UnknownAction", func(t *testing.T) {
		body := []any{2.0, "uuid-123", "Unknown", map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_InvalidUUID", func(t *testing.T) {
		body := []any{2.0, 123, "Unknown", map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("Request", func(t *testing.T) {
		body := []any{2.0, "uuid-123", types.Heartbeat, map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		parsed, err := machine.parseRawMessage(raw)
		assert.NoError(t, err)
		assert.Equal(t, types.Request, parsed.kind)
		assert.NotNil(t, parsed.action)
		assert.Equal(t, types.Heartbeat, *parsed.action)
		assert.Equal(t, "uuid-123", parsed.uuid)
		assert.NotEmpty(t, parsed.payload)
	})

	t.Run("Confirmation", func(t *testing.T) {
		body := []any{3.0, "uuid-456", map[string]any{"currentTime": "2025-07-24T12:34:56Z"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		parsed, err := machine.parseRawMessage(raw)
		assert.NoError(t, err)
		assert.Equal(t, types.Confirmation, parsed.kind)
		assert.Nil(t, parsed.action)
		assert.Equal(t, "uuid-456", parsed.uuid)
		assert.NotEmpty(t, parsed.payload)
	})
}

func TestHandleRequest(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("UnknownAction", func(t *testing.T) {
		err := machine.handleRequest(ctx, "Unknown", []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.handleRequest(ctx, types.BootNotification, []byte(`{"invalid": "payload"}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.handleRequest(ctx, types.BootNotification, []byte(`{
            "chargeBoxSerialNumber": "91234567",
            "chargePointModel": "Zappi",
            "chargePointSerialNumber": "91234567",
            "chargePointVendor": "Myenergi",
            "firmwareVersion": "5540",
            "iccid": "",
            "imsi": "",
            "meterType": "",
            "meterSerialNumber": "91234567"
        }`))
		assert.NoError(t, err)
	})
}

func TestHandleConfirmation(t *testing.T) {
	store := &mockStore{}
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(store),
	)
	ctx := context.Background()

	t.Run("NoMatchingUuid", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-unknown", []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-000", []byte(`{"invalid": "payload"}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-000", []byte(`{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval": 30,
			"status": "Accepted"
        }`))
		assert.NoError(t, err)
	})
}

func TestHandleHeartbeatRequest(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleHeartbeatRequest(ctx, []byte(`{}`))
		assert.NoError(t, err)
	})
}

func TestHandleHeartbeatConfirmation(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.HandleHeartbeatConfirmation(ctx, []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleHeartbeatConfirmation(ctx, []byte(`{"currentTime": "2025-07-24T12:34:56Z"}`))
		assert.NoError(t, err)
	})
}

func TestHandleBootNotificationRequest(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationRequest(ctx, []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationRequest(ctx, []byte(`{
			"chargeBoxSerialNumber": "91234567",
			"chargePointModel": "Zappi",
			"chargePointSerialNumber": "91234567",
			"chargePointVendor": "Myenergi",
			"firmwareVersion": "5540",
			"iccid": "",
			"imsi": "",
			"meterType": "",
			"meterSerialNumber": "91234567"
		}`))
		assert.NoError(t, err)
	})
}

func TestHandleBootNotificationConfirmation(t *testing.T) {
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithStore(&mockStore{}),
	)
	ctx := context.Background()

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationConfirmation(ctx, []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationConfirmation(ctx, []byte(`{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval": 30,
			"status": "Accepted"
		}`))
		assert.NoError(t, err)
	})
}
