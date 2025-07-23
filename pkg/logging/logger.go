package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

var levelVar slog.LevelVar

const (
	LogEnvDevelopment = "development"
	LogEnvProduction  = "production"
)

func SetupLogger(logLevel LogLevel, env string) {
	level := slog.LevelInfo

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	levelVar.Set(level)

	var handler slog.Handler
	if env == LogEnvDevelopment {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     &levelVar,
			AddSource: false,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     &levelVar,
			AddSource: false,
		})
	}

	slog.SetDefault(slog.New(NewSourceTrimmingHandler(handler)))

}

type SourceTrimmingHandler struct {
	slog.Handler
	trimPrefix string
}

func NewSourceTrimmingHandler(h slog.Handler) *SourceTrimmingHandler {
	wd, err := os.Getwd()
	if err != nil {
		wd = "" // fallback to no trim
	}
	return &SourceTrimmingHandler{
		Handler:    h,
		trimPrefix: wd + string(os.PathSeparator),
	}
}

func (h *SourceTrimmingHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := frames.Next()

		file := frame.File
		if file != "" && strings.HasPrefix(file, h.trimPrefix) {
			file = strings.TrimPrefix(file, h.trimPrefix)
		}

		if rel := strings.TrimPrefix(file, h.trimPrefix); rel != "" {
			file = rel
		}

		r.AddAttrs(slog.String("source", fmt.Sprintf("%s:%d", file, frame.Line)))
	}

	return h.Handler.Handle(ctx, r)
}
