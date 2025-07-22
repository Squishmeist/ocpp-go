package ocpp

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/sirupsen/logrus"
)

func detectOcppMessage(payload []byte) OcppMessage {
    if conf, err := unmarshalAndValidate[core.HeartbeatConfirmation](payload); err == nil {
        return OcppMessage{Type: HeartbeatConfirmation, Data: conf}
    }

    if req, err := unmarshalAndValidate[core.HeartbeatRequest](payload); err == nil {
        if string(payload) == "{}" || len(payload) == 0 {
            return OcppMessage{Type: HeartbeatRequest, Data: req}
        }
    }

    if bootReq, err := unmarshalAndValidate[core.BootNotificationRequest](payload); err == nil {
        return OcppMessage{Type: BootNotificationRequest, Data: bootReq}
    }

    if bootConf, err := unmarshalAndValidate[core.BootNotificationConfirmation](payload); err == nil {
        return OcppMessage{Type: BootNotificationConfirmation, Data: bootConf}
    }

    return OcppMessage{Type: Unknown, Data: nil}
}

func handleMessage(body Body, log *logrus.Logger) error {
	msg := detectOcppMessage([]byte(body.Payload))

	switch msg.Type {
		case HeartbeatConfirmation:
			conf, ok := asType[core.HeartbeatConfirmation](msg.Data)
			if !ok {
				log.Warn("Invalid HeartbeatConfirmation message")
				return nil
			}
			log.Infof("HeartbeatConfirmation: %v", conf.CurrentTime)
		case HeartbeatRequest:
			req, ok := asType[core.HeartbeatRequest](msg.Data)
			if !ok {
				log.Warn("Invalid HeartbeatRequest message")
				return nil
			}
			log.Infof("HeartbeatRequest: %v", req)
		case BootNotificationRequest:
			req, ok := asType[core.BootNotificationRequest](msg.Data)
			if !ok {
				log.Warn("Invalid BootNotificationRequest message")
				return nil
			}
			log.Infof("BootNotificationRequest: %v", req.ChargePointModel)
		case BootNotificationConfirmation:
			conf, ok := asType[core.BootNotificationConfirmation](msg.Data)
			if !ok {
				log.Warn("Invalid BootNotificationConfirmation message")
				return nil
			}
			log.Infof("BootNotificationConfirmation: %v", conf.CurrentTime)
		default:
			log.Warn("Unknown OCPP message type")
	}

	return nil

}

