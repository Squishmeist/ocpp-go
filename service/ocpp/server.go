package ocpp

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/squishmeist/ocpp-go/internal/core"

	"github.com/lorenzodonini/ocpp-go/ocppj"
)


func NewServer(log *logrus.Logger) *core.HttpServer {
	ocppj.SetLogger(log.WithField("logger", "ocppj"))

	server := core.NewHttpServer(
		func(s *core.HttpServer) {
			s.ServiceName = "ocpp-message-handler"
		},
	)
	topic(server, log)

	return server

}

func topic(server *core.HttpServer, log *logrus.Logger) {
	server.AddRoute(http.MethodPost, "/topic", func(ctx echo.Context) error {
		log.Info("Received request")
		var body Body
		if err := ctx.Bind(&body); err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid request")
		}

		err := handleMessage(body, log)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, "Failed to process request")
		}

		return ctx.String(http.StatusOK, "Triggered")
	})
}
