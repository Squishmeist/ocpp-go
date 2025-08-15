package message

import (
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	messagepb "github.com/squishmeist/ocpp-go/pkg/api/proto/message/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewServer(config utils.Configuration, client *core.AzureServiceBusClient) *core.GrpcServer {
	server := core.NewGrpcServer(
		core.WithGrpcServiceName("message"),
		core.WithGrpcPort(config.HttpServer.Port),
	)

	handler := NewMessageService(
		WithMessageClient(client),
		WithMessageInboundName(config.AzureServiceBus.TopicInbound.Name),
	)
	grpcTransport := NewMessageGrpcTransport(handler)
	messagepb.RegisterOCPPMessageServer(server.Grpc, grpcTransport)

	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	return server
}
