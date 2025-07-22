package ocpp

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)


func TestUnmarshalAndValidate_HeartbeatConfirmation(t *testing.T) {
    got, err := unmarshalAndValidate[core.HeartbeatConfirmation](
        []byte(`{"currentTime":"2025-07-22T12:34:56Z"}`),
    )
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

	want := "2025-07-22 12:34:56 +0000 UTC"
    if got.CurrentTime.String() != want {
        t.Errorf("unexpected currentTime: got %v, want %v", got.CurrentTime, want)
    }
}

func TestUnmarshalAndValidate_HeartbeatConfirmation2(t *testing.T) {
    _, err := unmarshalAndValidate[core.HeartbeatConfirmation](
        []byte(`{"unknown":"unknown"}`),
    )
	if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestUnmarshalAndValidate_HeartbeatConfirmation3(t *testing.T) {
    _, err := unmarshalAndValidate[core.HeartbeatConfirmation](
        []byte(`{"currentTime":"unknown"}`),
    )
	if err == nil {
        t.Errorf("expected error, got nil")
    }
}

func TestAsType_HeartbeatConfirmation(t *testing.T) {
    got, ok := asType[core.HeartbeatConfirmation](&core.HeartbeatConfirmation{
        CurrentTime: types.NewDateTime(time.Date(2025, 7, 22, 12, 34, 56, 0, time.UTC)),
    })
	if !ok {
		t.Error("Failed to cast to HeartbeatConfirmation")
	}

	want := "2025-07-22 12:34:56 +0000 UTC"
    if got.CurrentTime.String() != want {
        t.Errorf("unexpected currentTime: got %v, want %v", got.CurrentTime, want)
    }
}

func TestAsType_HeartbeatConfirmation2(t *testing.T) {
    // Pass a pointer to a value of a different type to test failed type assertion
    _, ok := asType[core.HeartbeatConfirmation](new(string))
    if ok {
        t.Error("Expected failure casting to HeartbeatConfirmation, but got success")
    }
}

