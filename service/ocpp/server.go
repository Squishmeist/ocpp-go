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
	test(server, state)

	return server
}

func test(server *core.HttpServer, state *State) {
	server.AddRoute(http.MethodPost, "/test", func(ctx echo.Context) error {
		body, err := deconstructBody(ctx)
		if err != nil {
			fmt.Println("Failed to deconstruct request body:", err)
			return ctx.String(http.StatusBadRequest, "Invalid request body")
		}

		switch body := body.(type) {
		case RequestBody:
			err := handleRequestBody(body, state)
			if err != nil {
				fmt.Println("Error handling RequestBody:", err)
				return ctx.String(http.StatusInternalServerError, "Error handling RequestBody")
			}
		case ConfirmationBody:
			err := handleConfirmationBody(body, state)
			if err != nil {
				fmt.Println("Error handling ConfirmationBody:", err)
				return ctx.String(http.StatusInternalServerError, "Failed to process")
			}
		default:
			return ctx.String(http.StatusBadRequest, "Unknown body type")
		}

		fmt.Println("State after processing: ", *state)
		return ctx.String(http.StatusOK, "Processed successfully")
	})
}

