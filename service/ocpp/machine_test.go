package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/squishmeist/ocpp-go/service/ocpp/types"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestHandleMessage(t *testing.T) {
	ctx, machine := setupMachineTest(t)
	meta := types.Meta{
		Id:           "test-id",
		Serialnumber: "test-serial",
	}

	err := machine.cache.AddRequest(ctx, meta, types.RequestBody{
		Uuid:   "uuid-456",
		Action: types.BootNotification,
		Payload: []byte(`{
			"chargeBoxSerialNumber": "91234567",
			"chargePointModel": "Zappi",
			"chargePointSerialNumber": "91234567",
			"chargePointVendor": "Myenergi",
			"firmwareVersion": "5540",
			"iccid": "",
			"imsi": "",
			"meterType": "",
			"meterSerialNumber": "91234567"
		}`),
	})
	assert.NoError(t, err)

	t.Run("InvalidKind", func(t *testing.T) {
		body := []any{"2.0", "uuid-000", types.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("NoKind", func(t *testing.T) {
		body := []any{"uuid-00", types.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_UnknownAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-000", types.Heartbeat}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_MissingPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-000", types.BootNotification, map[string]any{
			"chargeBoxSerialNumber":   "91234567",
			"chargePointModel":        "Zappi",
			"chargePointSerialNumber": "91234567",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
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
		err = machine.HandleMessage(ctx, meta, raw)
		assert.NoError(t, err)
	})

	t.Run("InvalidConfirmation_NoPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-000"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidConfirmation_MissingPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-000", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
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
		err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("Confirmation", func(t *testing.T) {
		body := []any{3.0, "uuid-456", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
			"status":      "Accepted",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		err = machine.HandleMessage(ctx, meta, raw)
		assert.NoError(t, err)
	})

}

func TestParseRawMessage(t *testing.T) {
	_, machine := setupMachineTest(t)

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
		body := []any{"not-a-number", "uuid-000"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("UnknownKind", func(t *testing.T) {
		body := []any{99.0, "uuid-000"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", 2, map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_UnknownAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", "Unknown", map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.parseRawMessage(raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_InvalidUUID", func(t *testing.T) {
		body := []any{2.0, 000, "Unknown", map[string]any{"custom": "value"}}
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
	ctx, machine := setupMachineTest(t)

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
	ctx, machine := setupMachineTest(t)

	err := machine.cache.AddRequest(ctx, types.Meta{
		Id:           "test-id",
		Serialnumber: "test-serial",
	}, types.RequestBody{
		Uuid:   "uuid-123",
		Action: types.BootNotification,
		Payload: []byte(`{
			"chargeBoxSerialNumber": "91234567",
			"chargePointModel": "Zappi",
			"chargePointSerialNumber": "91234567",
			"chargePointVendor": "Myenergi",
			"firmwareVersion": "5540",
			"iccid": "",
			"imsi": "",
			"meterType": "",
			"meterSerialNumber": "91234567"
		}`),
	})
	assert.NoError(t, err)

	t.Run("NoMatchingUuid", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-unknown", []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-000", []byte(`{"invalid": "payload"}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.handleConfirmation(ctx, "uuid-123", []byte(`{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval": 30,
			"status": "Accepted"
        }`))
		assert.NoError(t, err)
	})
}

func TestHandleHeartbeatRequest(t *testing.T) {
	ctx, machine := setupMachineTest(t)

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleHeartbeatRequest(ctx, []byte(`{}`))
		assert.NoError(t, err)
	})
}

func TestHandleHeartbeatConfirmation(t *testing.T) {
	ctx, machine := setupMachineTest(t)

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
	ctx := context.Background()
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithCache(&mockCache{}),
		WithStore(&mockStore{}),
	)

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
	ctx := context.Background()
	request := types.RequestBody{
		Uuid:   "uuid-000",
		Action: types.BootNotification,
		Payload: []byte(`{
			"chargeBoxSerialNumber": "91234567",
			"chargePointModel": "Zappi",
			"chargePointSerialNumber": "91234567",
			"chargePointVendor": "Myenergi",
			"firmwareVersion": "5540",
			"iccid": "",
			"imsi": "",
			"meterType": "",
			"meterSerialNumber": "91234567"
		}`),
	}
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithCache(&mockCache{}),
		WithStore(&mockStore{}),
	)

	t.Run("InvalidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationConfirmation(ctx, request, []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		err := machine.HandleBootNotificationConfirmation(ctx, request, []byte(`{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval": 30,
			"status": "Accepted"
		}`))
		assert.NoError(t, err)
	})
}

type mockCache struct {
	processed []string
	requests  map[string]types.RequestBody
}

func (m *mockCache) HasProcessed(ctx context.Context, id string) (bool, error) {
	if slices.Contains(m.processed, id) {
		return true, nil
	}

	return false, nil
}

func (m *mockCache) GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, error) {

	for _, request := range m.requests {
		if request.Uuid == uuid {
			return request, nil
		}
	}

	return types.RequestBody{}, fmt.Errorf("request not found")

}

func (m *mockCache) AddRequest(ctx context.Context, meta types.Meta, request types.RequestBody) error {
	if m.requests == nil {
		m.requests = make(map[string]types.RequestBody)
	}
	m.requests[request.Uuid] = request
	return nil
}

func (m *mockCache) RemoveRequest(ctx context.Context, meta types.Meta, confirmation types.ConfirmationBody) error {
	if m.requests == nil {
		return fmt.Errorf("no requests found")
	}
	delete(m.requests, confirmation.Uuid)
	return nil
}

type mockStore struct {
}

func (m *mockStore) AddChargepoint(ctx context.Context, request types.BootNotificationRequest) error {
	return nil
}

func setupMachineTest(t *testing.T) (context.Context, *OcppMachine) {
	ctx := context.Background()

	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithCache(&mockCache{}),
		WithStore(&mockStore{}),
	)
	return ctx, machine
}
