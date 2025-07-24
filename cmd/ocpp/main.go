package main

import (
	"context"
	"log/slog"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/pkg/logging"
	"github.com/squishmeist/ocpp-go/service/ocpp"
)

const (
	topicName        = "topic.1"
	subscriptionName = "subscription.1"
	endpoint         = "localhost:4318"
	serviceName      = "ocpp-machine"
	namespace        = "ocpp"
	connectionString = "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true;"
)

func main() {
	logging.SetupLogger(logging.LevelDebug, logging.LogEnvDevelopment)
	ctx := context.Background()

	t := core.NewTelemeter(serviceName, endpoint, namespace)
	tp := t.NewTracerProvider()
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown tracer provider", "error", err)
		}
	}()

	err := ocpp.Start(ctx, topicName, subscriptionName, connectionString, tp)
	if err != nil {
		panic(err)
	}
}
