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
	namespace        = "ocpp-go"
	serviceName      = "OCPPService"
	connectionString = "Endpoint=sb://localhost;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=SAS_KEY_VALUE;UseDevelopmentEmulator=true;"
)


func main() {
	logging.SetupLogger(logging.LevelDebug, logging.LogEnvDevelopment)
	
	t := core.NewTelemeter(serviceName, endpoint, namespace)
	tp := t.NewTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracer provider", "error", err)
		}
	}()

	state := &ocpp.State{}

	err := ocpp.Start(state, topicName, subscriptionName, connectionString, tp)
	if err != nil {
		panic(err)
	}
}

