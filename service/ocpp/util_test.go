package ocpp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

func TestDeconstructBody(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		wantType any
		wantErr  bool
	}{
		{
			name: "RequestBody",
			input: []any{
				2,
				"uuid-123",
				"Heartbeat",
				map[string]any{},
			},
			wantType: RequestBody{},
			wantErr:  false,
		},
		{
			name: "ConfirmationBody",
			input: []any{
				3,
				"uuid-456",
				map[string]any{"currentTime": "2025-07-22T11:25:25.230Z"},
			},
			wantType: ConfirmationBody{},
			wantErr:  false,
		},
		{
			name: "Invalid type",
			input: []any{
				99,
				"uuid-789",
			},
			wantType: nil,
			wantErr:  true,
		},
		{
			name: "Too few elements",
			input: []any{
				2,
			},
			wantType: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			bodyBytes, _ := json.Marshal(tt.input)
			got, err := deconstructBody(ctx, bodyBytes)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch tt.wantType.(type) {
			case RequestBody:
				rb, ok := got.(RequestBody)
				if !ok {
					t.Errorf("expected RequestBody, got %T", got)
				}
				// Check that Payload is []byte
				if rb.Payload == nil || len(rb.Payload) == 0 {
					t.Errorf("expected non-empty Payload []byte, got %v", rb.Payload)
				}
			case ConfirmationBody:
				cb, ok := got.(ConfirmationBody)
				if !ok {
					t.Errorf("expected ConfirmationBody, got %T", got)
				}
				// Check that Payload is []byte
				if cb.Payload == nil || len(cb.Payload) == 0 {
					t.Errorf("expected non-empty Payload []byte, got %v", cb.Payload)
				}
			}
		})
	}
}

func TestUnmarshalAndValidate_HeartbeatConfirmation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "valid HeartbeatConfirmation",
			payload: `{"currentTime":"2025-07-22T12:34:56Z"}`,
			wantErr: false,
		},
		{
			name:    "missing currentTime",
			payload: `{}`,
			wantErr: true,
		},
		{
			name:    "invalid currentTime format",
			payload: `{"currentTime":"not-a-date"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := unmarshalAndValidate[core.HeartbeatConfirmation]([]byte(tt.payload))
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestUnmarshalAndValidate_BootNotificationConfirmation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "valid BootNotificationConfirmation",
			payload: `{"status":"Accepted","currentTime":"2025-07-22T12:34:56Z"}`,
			wantErr: false,
		},
		{
			name:    "missing status",
			payload: `{"currentTime":"2025-07-22T12:34:56Z"}`,
			wantErr: true,
		},
		{
			name:    "invalid currentTime format",
			payload: `{"status":"Accepted","currentTime":"not-a-date"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := unmarshalAndValidate[core.BootNotificationConfirmation]([]byte(tt.payload))
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
