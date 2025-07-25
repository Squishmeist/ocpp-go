package utils_test

import (
	"os"
	"testing"

	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	logging.SetupLogger(logging.LevelInfo, logging.LogEnvProduction)

	t.Run("Returns config for valid config file", func(t *testing.T) {
		// Act create a config file
		file, err := os.Create("./example.yaml")
		assert.NoError(t, err)

		_, err = file.WriteString("HTTP_SERVER:\n  PORT: \":8080\"\n")
		assert.NoError(t, err)

		// Act
		config := utils.GetConfig(".", "example", "yaml")

		// Assert
		assert.Equal(t, ":8080", config.HttpServer.Port)

		// Cleanup
		err = os.Remove("./example.yaml")
		assert.NoError(t, err)
	})

	t.Run("Returns config for valid azure-service-bus file", func(t *testing.T) { // Act create a config file
		file, err := os.Create("./example.yaml")
		assert.NoError(t, err)

		_, err = file.WriteString("AZURE_SERVICE_BUS:\n  CONNECTION_STRING: \"Endpoint=sb://localhost;\"\n  TOPIC_INBOUND:\n    NAME: \"topic.inbound\"\n    SUBSCRIPTION: \"subscription.inbound\"\n  TOPIC_OUTBOUND:\n    NAME: \"topic.outbound\"\n    SUBSCRIPTION: \"subscription.outbound\"\n")
		assert.NoError(t, err)

		// Act
		config := utils.GetConfig(".", "example", "yaml")

		// Assert
		assert.Equal(t, "Endpoint=sb://localhost;", config.AzureServiceBus.ConnectionString)

		// Cleanup
		err = os.Remove("./example.yaml")
		assert.NoError(t, err)
	})

}
