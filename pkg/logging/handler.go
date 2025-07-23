package logging

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

type SloggerControlRequest struct {
	Level LogLevel `json:"level"`
}

type SloggerControlResponse struct {
	Message string `json:"message"`
}

func Handler(c echo.Context) error {
	req := new(SloggerControlRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, SloggerControlResponse{Message: "Invalid request"})
	}

	SetupLogger(req.Level, LogEnvProduction)

	return c.JSON(200, SloggerControlResponse{Message: fmt.Sprintf("Log level set to %s", req.Level)})
}
