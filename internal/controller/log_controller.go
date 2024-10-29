package controller

import (
	"fmt"
	"open_url_service/pkg/logger"
)

type logController struct {
}

func (l logController) Serve() error {
	var (
		lf = logger.NewFields(logger.EventName("logProcessor"))
	)
	//var payload broker.MessagePayload
	//_ = json.Unmarshal(data.Body, &payload)

	logger.Info(fmt.Sprintf("Payload Data %+v)"), lf...)

	return nil
}

//func NewLogController() contract.MessageController {
//	return &logController{}
//}
