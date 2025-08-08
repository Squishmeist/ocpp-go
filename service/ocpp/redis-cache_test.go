package ocpp

import (
	"context"
	"testing"

	v16 "github.com/squishmeist/ocpp-go/service/ocpp/v1.6"
	"github.com/squishmeist/ocpp-go/service/ocpp/v1.6/core"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestRedisCache(t *testing.T) {
	ctx, redis := setupRedisTest(t)

	// Add a request
	err := redis.AddRequest(ctx, v16.Meta{
		Id:           "test-id",
		Serialnumber: "test-serialnumber",
	}, v16.RequestBody{
		Uuid:    "test-uuid",
		Action:  core.Heartbeat,
		Payload: []byte("{}"),
	})
	assert.NoError(t, err, "expected no error when adding request")

	t.Cleanup(func() {
		keys, err := redis.client.Keys(ctx, "test-id*").Result()
		if err == nil && len(keys) > 0 {
			redis.client.Del(ctx, keys...)
		}

		keys, err = redis.client.Keys(ctx, "request:test-uuid*").Result()
		if err == nil && len(keys) > 0 {
			redis.client.Del(ctx, keys...)
		}

		if err := redis.Close(); err != nil {
			t.Logf("Warning: failed to close Redis connection: %v", err)
		}
	})

	t.Run("HasProcessed_Unknown Id", func(t *testing.T) {
		processed, err := redis.HasProcessed(ctx, "unknown-id")
		assert.Nil(t, err, "expected no error for unknown id")
		assert.False(t, processed, "expected processed to be false for unknown id")
	})

	t.Run("HasProcessed_Known Id", func(t *testing.T) {
		processed, err := redis.HasProcessed(ctx, "test-id")
		assert.Nil(t, err, "expected no error for valid id")
		assert.True(t, processed, "expected processed to be true for valid id")
	})

	t.Run("AddRequest_Valid", func(t *testing.T) {
		err := redis.AddRequest(ctx, v16.Meta{
			Id:           "test-id2",
			Serialnumber: "test-serialnumber",
		}, v16.RequestBody{
			Uuid:    "test-uuid2",
			Action:  core.Heartbeat,
			Payload: []byte("{}"),
		})
		assert.NoError(t, err, "expected no error for valid input")
	})

	t.Run("RemoveRequest_Valid", func(t *testing.T) {
		err := redis.AddRequest(ctx, v16.Meta{
			Id:           "test-id3",
			Serialnumber: "test-serialnumber",
		}, v16.RequestBody{
			Uuid:    "test-uuid3",
			Action:  core.Heartbeat,
			Payload: []byte("{}"),
		})
		assert.NoError(t, err, "expected no error for valid input")

		err = redis.RemoveRequest(ctx, v16.Meta{
			Id:           "test-id3",
			Serialnumber: "test-serialnumber",
		}, v16.ConfirmationBody{
			Uuid: "test-uuid3",
		})
		assert.NoError(t, err, "expected no error for valid input")
	})

	t.Run("RemoveRequest_Unknown", func(t *testing.T) {
		err := redis.RemoveRequest(ctx, v16.Meta{
			Id:           "test-unknown",
			Serialnumber: "test-serialnumber",
		}, v16.ConfirmationBody{
			Uuid: "unknown-uuid",
		})
		assert.Error(t, err, "expected error for unknown uuid")
	})

	t.Run("GetRequestFromUuid_Valid", func(t *testing.T) {
		request, err := redis.GetRequestFromUuid(ctx, "test-uuid")
		assert.NoError(t, err, "expected no error for valid input")
		assert.Equal(t, "test-uuid", request.Uuid)
		assert.Equal(t, v16.ActionKind(core.Heartbeat), request.Action)
		assert.Equal(t, []byte("{}"), request.Payload)
	})

	t.Run("GetRequestFromUuid_Unknown", func(t *testing.T) {
		_, err := redis.GetRequestFromUuid(ctx, "unknown-uuid")
		assert.Error(t, err, "expected error for unknown uuid")
	})
}

func setupRedisTest(t *testing.T) (context.Context, *RedisCache) {
	ctx := context.Background()
	redis := NewRedisCache(noop.NewTracerProvider(), "localhost:6379")

	// Test connection
	err := redis.client.Ping(ctx).Err()
	assert.Nil(t, err, "expected successful connection to Redis")

	return ctx, redis
}
