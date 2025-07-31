package utils

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Telemetry       TelemetryConfiguration
	AzureServiceBus AzureServiceBusConfiguration
	HttpServer      HttpServer
	Database        DatabaseConfiguration
}

type Topic struct {
	Name         string
	Subscription string
}

type TelemetryConfiguration struct {
	ENDPOINT string
}

type AzureServiceBusConfiguration struct {
	ConnectionString string
	TopicInbound     Topic
	TopicOutbound    Topic
}

type HttpServer struct {
	Port string
	Host string
}

type DatabaseConfiguration struct {
	Address  string
	Protocol string
	Driver   string
	PoolSize int
}

func initiateConfigDefaults(configName string, configPath []string, configType string) *viper.Viper {

	viperObj := viper.New()

	for _, path := range configPath {
		viperObj.AddConfigPath(path)
		contents, err := os.ReadDir(path)

		if err != nil {
			slog.Error("Error reading directory", "error", err, "path", path)
		}

		for _, content := range contents {
			if content.IsDir() {
				continue
			}

			if content.Name() == fmt.Sprintf("%s.%s", configName, configType) {
				slog.Info("Found config file at path", "path", path, "file", content.Name())
			}
		}
	}

	viperObj.SetConfigName(configName)
	viperObj.SetConfigType(configType)
	viperObj.SetDefault("port", 8000)

	return viperObj
}

func GetConfig(path, name, extn string) Configuration {

	paths := []string{path, "/etc/conf"}

	viperObj := initiateConfigDefaults(name, paths, extn)

	err := viperObj.ReadInConfig()
	if err != nil {
		slog.Error("Configuration file not found in path location, using default values", "error", err, "path", path, "name", name, "extn", extn)
	}

	return Configuration{
		Telemetry: TelemetryConfiguration{
			ENDPOINT: viperObj.GetString("TELEMETRY.ENDPOINT"),
		},
		AzureServiceBus: AzureServiceBusConfiguration{
			ConnectionString: viperObj.GetString("AZURE_SERVICE_BUS.CONNECTION_STRING"),
			TopicInbound: Topic{
				Name:         viperObj.GetString("AZURE_SERVICE_BUS.TOPIC_INBOUND.NAME"),
				Subscription: viperObj.GetString("AZURE_SERVICE_BUS.TOPIC_INBOUND.SUBSCRIPTION"),
			},
			TopicOutbound: Topic{
				Name:         viperObj.GetString("AZURE_SERVICE_BUS.TOPIC_OUTBOUND.NAME"),
				Subscription: viperObj.GetString("AZURE_SERVICE_BUS.TOPIC_OUTBOUND.SUBSCRIPTION"),
			},
		},
		HttpServer: HttpServer{
			Port: viperObj.GetString("HTTP_SERVER.PORT"),
			Host: viperObj.GetString("HTTP_SERVER.HOST"),
		},
		Database: DatabaseConfiguration{
			Address:  viperObj.GetString("DATABASE.ADDR"),
			Protocol: viperObj.GetString("DATABASE.PROTOCOL"),
			Driver:   viperObj.GetString("DATABASE.DRIVER"),
			PoolSize: viperObj.GetInt("DATABASE.POOL_SIZE"),
		},
	}
}
