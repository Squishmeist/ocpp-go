package main

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/squishmeist/ocpp-go/service/ocpp"
)

const (
	defaultListenPort          = 8887
	defaultHeartbeatInterval   = 600
	envVarServerPort           = "SERVER_LISTEN_PORT"
	envVarTls                  = "TLS_ENABLED"
	envVarCaCertificate        = "CA_CERTIFICATE_PATH"
	envVarServerCertificate    = "SERVER_CERTIFICATE_PATH"
	envVarServerCertificateKey = "SERVER_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger

func main() {
	server := ocpp.NewServer(log)
	log.Infof("starting server on port %v", defaultListenPort)
	server.Start(":" + strconv.Itoa(defaultListenPort))
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.DebugLevel)
}
