package ocpp

import (
	"encoding/json"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/sirupsen/logrus"
)

func asType[T any](data interface{}) (*T, bool) {
    val, ok := data.(*T)
    return val, ok
}

func tryUnmarshal[T any](payload []byte) (*T, error) {
    var obj T
    err := json.Unmarshal(payload, &obj)
    if err != nil {
        return nil, err
    }
    return &obj, nil
}

func detectOcppMessage(payload []byte) OcppMessage {
    if conf, err := tryUnmarshal[core.HeartbeatConfirmation](payload); err == nil && conf.CurrentTime != nil {
        return OcppMessage{Type: HeartbeatConfirmation, Data: conf}
    }

    if req, err := tryUnmarshal[core.HeartbeatRequest](payload); err == nil {
        if string(payload) == "{}" || len(payload) == 0 {
            return OcppMessage{Type: HeartbeatRequest, Data: req}
        }
    }

    if bootReq, err := tryUnmarshal[core.BootNotificationRequest](payload); err == nil && bootReq.ChargePointModel != "" {
        return OcppMessage{Type: BootNotificationRequest, Data: bootReq}
    }

    if bootConf, err := tryUnmarshal[core.BootNotificationConfirmation](payload); err == nil && bootConf.CurrentTime != nil {
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

