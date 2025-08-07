package ocpp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/service/ocpp/types"
	"go.opentelemetry.io/otel/trace"
)

type RedisCache struct {
	Tracer trace.Tracer
	client *redis.Client
}

func NewRedisCache(tp trace.TracerProvider, address string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisCache{
		Tracer: tp.Tracer("cache"),
		client: rdb,
	}
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) HasProcessed(ctx context.Context, id string) (bool, error) {
	ctx, span := core.TraceCache(ctx, c.Tracer, "Cache.HasProcessed")
	defer span.End()

	val, err := c.client.Get(ctx, id).Result()

	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return val == "1", nil
}

func (c *RedisCache) addProcessed(ctx context.Context, id string) error {
	ctx, span := core.TraceCache(ctx, c.Tracer, "Cache.AddProcessed")
	defer span.End()
	if err := c.client.Set(ctx, id, "1", 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("error trying to set message with id %s: %w", id, err)
	}
	return nil
}

func (c *RedisCache) GetRequestFromUuid(ctx context.Context, uuid string) (types.RequestBody, error) {
	ctx, span := core.TraceCache(ctx, c.Tracer, "Cache.GetRequestFromUuid")
	defer span.End()

	result, err := c.client.HGetAll(ctx, "request:"+uuid).Result()

	if err != nil {
		return types.RequestBody{}, err
	}

	if len(result) == 0 {
		return types.RequestBody{}, fmt.Errorf("request not found")
	}

	if _, ok := result["uuid"]; !ok {
		return types.RequestBody{}, fmt.Errorf("uuid not found in request data")
	}
	if _, ok := result["action"]; !ok {
		return types.RequestBody{}, fmt.Errorf("action not found in request data")
	}
	if _, ok := result["payload"]; !ok {
		return types.RequestBody{}, fmt.Errorf("payload not found in request data")
	}

	return types.RequestBody{
		Uuid:    uuid,
		Action:  types.ActionKind(result["action"]),
		Payload: []byte(result["payload"]),
	}, nil
}

func (c *RedisCache) AddRequest(ctx context.Context, meta types.Meta, request types.RequestBody) error {
	ctx, span := core.TraceCache(ctx, c.Tracer, "Cache.AddRequest")
	defer span.End()

	requestMap := map[string]any{
		"uuid":    request.Uuid,
		"action":  string(request.Action),
		"payload": string(request.Payload),
	}

	if err := c.client.HSet(ctx, "request:"+request.Uuid, requestMap).Err(); err != nil {
		return err
	}
	if err := c.client.Expire(ctx, "request:"+request.Uuid, 24*time.Hour).Err(); err != nil {
		return err
	}

	if err := c.addProcessed(ctx, meta.Id); err != nil {
		return fmt.Errorf("error adding message to processed cache: %w", err)
	}

	return nil
}

func (c *RedisCache) RemoveRequest(ctx context.Context, meta types.Meta, request types.ConfirmationBody) error {
	ctx, span := core.TraceCache(ctx, c.Tracer, "Cache.RemoveRequest")
	defer span.End()

	exists, err := c.client.Exists(ctx, "request:"+request.Uuid).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("request not found")
	}

	if err := c.client.Del(ctx, "request:"+request.Uuid).Err(); err != nil {
		return err
	}
	if err := c.addProcessed(ctx, meta.Id); err != nil {
		return fmt.Errorf("error adding message to processed cache: %w", err)
	}

	return nil
}
