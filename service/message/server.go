package message

import (
	"github.com/squishmeist/ocpp-go/internal/core"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	ocpppb "github.com/squishmeist/ocpp-go/pkg/api/proto/ocpp/v1"
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
	ocpppb.RegisterOCPPMessageServer(server.Grpc, grpcTransport)

	return server
}
