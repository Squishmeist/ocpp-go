package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/pkg/logging"
	"github.com/squishmeist/ocpp-go/service/ocpp"
)

func main() {
	logging.SetupLogger(logging.LevelDebug, logging.LogEnvDevelopment)
	ctx := context.Background()

	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "ocpp"
	}
	conf := utils.GetConfig("./config", configName, "yaml")

	t := core.NewTelemeter("ocpp-machine", conf.Telemetry.ENDPOINT, "ocpp")
	tp := t.NewTracerProvider()
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown tracer provider", "error", err)
		}
	}()

	err := ocpp.Start(ctx, tp, conf)
	if err != nil {
		panic(err)
	}
}
