package main

import (
	"strconv"

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


func main() {
	server := ocpp.NewServer()
	server.Start(":" + strconv.Itoa(defaultListenPort))
}

