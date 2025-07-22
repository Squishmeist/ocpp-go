package client

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/logging"
	"github.com/sirupsen/logrus"

	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

const (
	envVarClientID             = "CLIENT_ID"
	envVarCentralSystemUrl     = "CENTRAL_SYSTEM_URL"
	envVarTls                  = "TLS_ENABLED"
	envVarCACertificate        = "CA_CERTIFICATE_PATH"
	envVarClientCertificate    = "CLIENT_CERTIFICATE_PATH"
	envVarClientCertificateKey = "CLIENT_CERTIFICATE_KEY_PATH"
)

var log *logrus.Logger

func setupChargePoint(chargePointID string) ocpp16.ChargePoint {
	return ocpp16.NewChargePoint(chargePointID, nil, nil)
}

func setupTlsChargePoint(chargePointID string) ocpp16.ChargePoint {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caPath, ok := os.LookupEnv(envVarCACertificate)
	if ok {
		caCert, err := os.ReadFile(caPath)
		if err != nil {
			log.Warn(err)
		} else if !certPool.AppendCertsFromPEM(caCert) {
			log.Info("no ca.cert file found, will use system CA certificates")
		}
	} else {
		log.Info("no ca.cert file found, will use system CA certificates")
	}
	// Load client certificate
	clientCertPath, ok1 := os.LookupEnv(envVarClientCertificate)
	clientKeyPath, ok2 := os.LookupEnv(envVarClientCertificateKey)
	var clientCertificates []tls.Certificate
	if ok1 && ok2 {
		certificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err == nil {
			clientCertificates = []tls.Certificate{certificate}
		} else {
			log.Infof("couldn't load client TLS certificate: %v", err)
		}
	}
	// Create client with TLS config
	client := ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: clientCertificates,
	})
	return ocpp16.NewChargePoint(chargePointID, nil, client)
}

// exampleRoutine simulates a charge point flow, where
func exampleRoutine(chargePoint ocpp16.ChargePoint, stateHandler *ChargePointHandler) {
	// Boot
	bootConf, err := chargePoint.BootNotification("model1", "vendor1")
	checkError(err)
	logDefault(bootConf.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Send log notification
	_, err = chargePoint.LogStatusNotification(logging.UploadLogStatusUploading, 1)
	checkError(err)
	// Wait for some time ...
	time.Sleep(5 * time.Minute)
}

// Start function
func main() {
	// Load config
	id, ok := os.LookupEnv(envVarClientID)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarClientID)
		return
	}
	csUrl, ok := os.LookupEnv(envVarCentralSystemUrl)
	if !ok {
		log.Printf("no %v environment variable found, exiting...", envVarCentralSystemUrl)
		return
	}
	// Check if TLS enabled
	t, _ := os.LookupEnv(envVarTls)
	tlsEnabled, _ := strconv.ParseBool(t)
	// Prepare OCPP 1.6 charge point (chargePoint variable is defined in handler.go)
	if tlsEnabled {
		chargePoint = setupTlsChargePoint(id)
	} else {
		chargePoint = setupChargePoint(id)
	}
	// Setup some basic state management
	connectors := map[int]*ConnectorInfo{
		1: {status: core.ChargePointStatusAvailable, availability: core.AvailabilityTypeOperative, currentTransaction: 0},
	}
	handler := &ChargePointHandler{
		status:               core.ChargePointStatusAvailable,
		connectors:           connectors,
		configuration:        getDefaultConfig(),
		errorCode:            core.NoError,
		localAuthList:        []localauth.AuthorizationData{},
		localAuthListVersion: 0}
	// Support callbacks for all OCPP 1.6 profiles
	chargePoint.SetCoreHandler(handler)
	chargePoint.SetFirmwareManagementHandler(handler)
	chargePoint.SetLocalAuthListHandler(handler)
	chargePoint.SetReservationHandler(handler)
	chargePoint.SetRemoteTriggerHandler(handler)
	chargePoint.SetSmartChargingHandler(handler)
	// OCPP 1.6j Security extension
	chargePoint.SetCertificateHandler(handler)
	chargePoint.SetLogHandler(handler)
	chargePoint.SetSecureFirmwareHandler(handler)
	chargePoint.SetExtendedTriggerMessageHandler(handler)
	chargePoint.SetSecurityHandler(handler)

	ocppj.SetLogger(log.WithField("logger", "ocppj"))
	ws.SetLogger(log.WithField("logger", "websocket"))
	// Connects to central system
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Errorln(err)
	} else {
		log.Infof("connected to central system at %v", csUrl)
		exampleRoutine(chargePoint, handler)
		// Disconnect
		chargePoint.Stop()
		log.Infof("disconnected from central system")
	}
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Set this to DebugLevel if you want to retrieve verbose logs from the ocppj and websocket layers
	log.SetLevel(logrus.DebugLevel)
}

// Utility functions
func logDefault(feature string) *logrus.Entry {
	return log.WithField("message", feature)
}
