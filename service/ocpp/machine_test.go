package ocpp

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	v16 "github.com/squishmeist/ocpp-go/service/ocpp/v1.6"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestHandleMessage(t *testing.T) {
	ctx, machine := setupMachineTest(t)
	meta := v16.Meta{
		Id:           "test-id",
		Serialnumber: "test-serial",
	}

	err := machine.cache.AddRequest(ctx, meta, v16.RequestBody{
		Uuid:   "uuid-456",
		Action: core.BootNotification,
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
		body := []any{"2.0", "uuid-000", core.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("NoKind", func(t *testing.T) {
		body := []any{"uuid-00", core.Heartbeat, map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_UnknownAction", func(t *testing.T) {
		body := []any{2.0, "uuid-000", map[string]any{}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_NoPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-000", core.Heartbeat}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidRequest_MissingPayload", func(t *testing.T) {
		body := []any{2.0, "uuid-000", core.BootNotification, map[string]any{
			"chargeBoxSerialNumber":   "91234567",
			"chargePointModel":        "Zappi",
			"chargePointSerialNumber": "91234567",
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("Request", func(t *testing.T) {
		body := []any{2.0, "uuid-123", core.BootNotification, map[string]any{
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
		bytes, err := machine.HandleMessage(ctx, meta, raw)
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("InvalidConfirmation_NoPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-000"}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	t.Run("InvalidConfirmation_MissingPayload", func(t *testing.T) {
		body := []any{3.0, "uuid-000", map[string]any{
			"currentTime": "2024-04-02T11:44:38Z",
			"interval":    30,
		}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		_, err = machine.HandleMessage(ctx, meta, raw)
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
		_, err = machine.HandleMessage(ctx, meta, raw)
		assert.Error(t, err)
	})

	// t.Run("Confirmation", func(t *testing.T) {
	// 	body := []any{3.0, "uuid-456", map[string]any{
	// 		"currentTime": "2024-04-02T11:44:38Z",
	// 		"interval":    30,
	// 		"status":      "Accepted",
	// 	}}
	// 	raw, err := json.Marshal(body)
	// 	assert.NoError(t, err)
	// 	bytes, err := machine.HandleMessage(ctx, meta, raw)
	// 	assert.NoError(t, err)
	// 	assert.NotEmpty(t, bytes)
	// })

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
		body := []any{2.0, "uuid-123", core.Heartbeat, map[string]any{"custom": "value"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		parsed, err := machine.parseRawMessage(raw)
		assert.NoError(t, err)
		assert.Equal(t, v16.Request, parsed.kind)
		assert.NotNil(t, parsed.action)
		assert.Equal(t, v16.ActionKind(core.Heartbeat), *parsed.action)
		assert.Equal(t, "uuid-123", parsed.uuid)
		assert.NotEmpty(t, parsed.payload)
	})

	t.Run("Confirmation", func(t *testing.T) {
		body := []any{3.0, "uuid-456", map[string]any{"currentTime": "2025-07-24T12:34:56Z"}}
		raw, err := json.Marshal(body)
		assert.NoError(t, err)
		parsed, err := machine.parseRawMessage(raw)
		assert.NoError(t, err)
		assert.Equal(t, v16.Confirmation, parsed.kind)
		assert.Nil(t, parsed.action)
		assert.Equal(t, "uuid-456", parsed.uuid)
		assert.NotEmpty(t, parsed.payload)
	})
}

func TestHandleRequest(t *testing.T) {
	ctx, machine := setupMachineTest(t)
	meta := v16.Meta{
		Id:           "test-id",
		Serialnumber: "test-serial",
	}

	t.Run("InvalidPayload", func(t *testing.T) {
		msg := parsedMessage{
			kind:    v16.Request,
			action:  v16.ActionKind(core.BootNotification).ToPtr(),
			uuid:    "uuid-000",
			payload: []byte(`{"invalid": "payload"}`),
		}
		_, err := machine.handleRequest(ctx, meta, msg)
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		msg := parsedMessage{
			kind:   v16.Request,
			action: v16.ActionKind(core.BootNotification).ToPtr(),
			uuid:   "uuid-000",
			payload: []byte(`{
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

		body, err := machine.handleRequest(ctx, meta, msg)
		assert.NoError(t, err)
		assert.NotEmpty(t, body)
	})
}

// func TestHandleConfirmation(t *testing.T) {
// 	ctx, machine := setupMachineTest(t)
// 	meta := v16.Meta{
// 		Id:           "test-id",
// 		Serialnumber: "test-serial",
// 	}

// 	err := machine.cache.AddRequest(ctx, v16.Meta{
// 		Id:           "test-id",
// 		Serialnumber: "test-serial",
// 	}, v16.RequestBody{
// 		Uuid:   "uuid-123",
// 		Action: core.BootNotification,
// 		Payload: []byte(`{
// 			"chargeBoxSerialNumber": "91234567",
// 			"chargePointModel": "Zappi",
// 			"chargePointSerialNumber": "91234567",
// 			"chargePointVendor": "Myenergi",
// 			"firmwareVersion": "5540",
// 			"iccid": "",
// 			"imsi": "",
// 			"meterType": "",
// 			"meterSerialNumber": "91234567"
// 		}`),
// 	})
// 	assert.NoError(t, err)

// 	t.Run("NoMatchingUuid", func(t *testing.T) {
// 		err := machine.handleConfirmation(ctx, meta, "uuid-unknown", []byte(`{}`))
// 		assert.Error(t, err)
// 	})

// 	t.Run("InvalidPayload", func(t *testing.T) {
// 		err := machine.handleConfirmation(ctx, meta, "uuid-000", []byte(`{"invalid": "payload"}`))
// 		assert.Error(t, err)
// 	})

// 	t.Run("ValidPayload", func(t *testing.T) {
// 		err := machine.handleConfirmation(ctx, meta, "uuid-123", []byte(`{
// 			"currentTime": "2024-04-02T11:44:38Z",
// 			"interval": 30,
// 			"status": "Accepted"
//         }`))
// 		assert.NoError(t, err)
// 	})
// }

func TestHandleHeartbeatRequest(t *testing.T) {
	ctx, machine := setupMachineTest(t)

	t.Run("ValidPayload", func(t *testing.T) {
		confirmation, err := machine.handleHeartbeatRequest(ctx, "test-serial", []byte(`{}`))
		assert.NoError(t, err)
		assert.NotEmpty(t, confirmation)
	})
}

func TestHandleBootNotificationRequest(t *testing.T) {
	ctx, machine := setupMachineTest(t)

	t.Run("InvalidPayload", func(t *testing.T) {
		_, err := machine.handleBootNotificationRequest(ctx, []byte(`{}`))
		assert.Error(t, err)
	})

	t.Run("ValidPayload", func(t *testing.T) {
		body, err := machine.handleBootNotificationRequest(ctx, []byte(`{
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
		assert.NotEmpty(t, body)
	})
}

type mockCache struct {
	processed []string
	requests  map[string]v16.RequestBody
}

func (m *mockCache) HasProcessed(ctx context.Context, id string) (bool, error) {
	if slices.Contains(m.processed, id) {
		return true, nil
	}

	return false, nil
}

func (m *mockCache) GetRequestFromUuid(ctx context.Context, uuid string) (v16.RequestBody, error) {

	for _, request := range m.requests {
		if request.Uuid == uuid {
			return request, nil
		}
	}

	return v16.RequestBody{}, fmt.Errorf("request not found")

}

func (m *mockCache) AddRequest(ctx context.Context, meta v16.Meta, request v16.RequestBody) error {
	if m.requests == nil {
		m.requests = make(map[string]v16.RequestBody)
	}
	m.requests[request.Uuid] = request
	return nil
}

func (m *mockCache) RemoveRequest(ctx context.Context, meta v16.Meta, confirmation v16.ConfirmationBody) error {
	if m.requests == nil {
		return fmt.Errorf("no requests found")
	}
	delete(m.requests, confirmation.Uuid)
	return nil
}

type mockStore struct {
}

func (m *mockStore) AddChargepoint(ctx context.Context, request core.BootNotificationRequest) error {
	return nil
}

func (m *mockStore) UpdateLastHeartbeat(ctx context.Context, serialnumber string, payload core.HeartbeatConfirmation) error {
	return nil
}

func setupMachineTest(t *testing.T) (context.Context, *OcppMachine) {
	ctx := context.Background()
	machine := NewOcppMachine(
		WithTracerProvider(noop.NewTracerProvider()),
		WithCache(&mockCache{}),
		WithStore(&mockStore{}),
	)
	assert.NotNil(t, machine)
	return ctx, machine
}
