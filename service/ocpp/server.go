package ocpp

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/squishmeist/ocpp-go/internal/core"
)


func NewServer() *core.HttpServer {
	state := &State{}

	server := core.NewHttpServer(
		func(s *core.HttpServer) {
			s.ServiceName = "ocpp-message-handler"
		},
	)
	topic(server, state)

	return server
}

func topic(server *core.HttpServer, state *State) {
	server.AddRoute(http.MethodPost, "/topic", func(ctx echo.Context) error {
		// Deconstruct the request body
		body, err := deconstructBody(ctx)
		if err != nil {
			fmt.Println("Failed to deconstruct request body:", err)
			return ctx.String(http.StatusBadRequest, "Invalid request body")
		}

		// Process the body based on its type
		switch body := body.(type) {
		case RequestBody:
			handleRequestBody(body, state)
		case ConfirmationBody:
			handleConfirmationBody(body, state)
		default:
			return ctx.String(http.StatusBadRequest, "Unknown body type")
		}

		fmt.Print("State after processing: ", *state)
		return ctx.String(http.StatusOK, "Processed successfully")
	})
}

