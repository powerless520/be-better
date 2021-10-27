package dana

import (
	"be-better/core/global"
	"be-better/utils/reflectUtil"
	"errors"
	"gitlab.ops.kingnet.com/danadata/go-sdk/com.kingnetdc.danadata.gosdk/model"
	"strconv"
	"time"
)

func Send(event, did, ouid string, properties map[string]interface{}) (int, error) {

	msg := model.NewMessage(did, event, ouid, "facm", time.Now().UnixNano()/1e6)

	if properties != nil {
		for k, v := range properties {
			msg.WithProperty(k, v)
		}
	}

	status, err := global.GlobalDANA.Facm.Send(msg)
	if err != nil {
		return 0, err
	}

	if !status {
		errMsg := "dana send fail, data: " + msg.ToString()
		global.GlobalLogger.Error(errMsg)
		return 0, errors.New(errMsg)
	}

	return 1, nil
}

func SendEvent(event, did, ouid string, key string, properties map[string]interface{}) (int, error) {

	msg := model.NewMessage(did, event, ouid, "facm", time.Now().UnixNano()/1e6)

	if properties != nil {
		for k, v := range properties {
			msg.WithProperty(k, v)
		}
	}

	var (
		status bool
		err    error
	)

	if key == "" {
		status, err = global.GlobalDANA.FacmEvent.Send(msg)
	} else {
		status, err = global.GlobalDANA.FacmEvent.SendWithKey(msg, key)
	}

	if err != nil {
		return 0, err
	}

	if !status {
		errMsg := "dana send fail, data: " + msg.ToString()
		global.GlobalLogger.Error(errMsg)
		return 0, errors.New(errMsg)
	}
	global.GlobalLogger.Info("send data success")
	return 1, nil
}

func PushIdcardAuthenticationRequest(idcardAuthenticationRequestEvent IdcardAuthenticationRequestEvent) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(idcardAuthenticationRequestEvent)
	if err != nil {
		return 0, err
	}
	return Send("IdcardAuthenticationRequest", "-1", "-1", properties)
}

func PushIdcardAuthentication(idcardAuthenticationEvent IdcardAuthenticationEvent) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(idcardAuthenticationEvent)
	if err != nil {
		return 0, err
	}
	return Send("IdcardAuthentication", "-1", "-1", properties)
}

func PushIdcardAuthenticationNotify(idcardAuthenticationNotifyEvent IdcardAuthenticationNotifyEvent) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(idcardAuthenticationNotifyEvent)
	if err != nil {
		return 0, err
	}
	return Send("IdcardAuthenticationNotify", "-1", "-1", properties)
}

func PushLoginOutEvent(loginoutEvent LoginOutEvent) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(loginoutEvent)
	if err != nil {
		return 0, err
	}

	return SendEvent("LoginOut", "-1", "-1", strconv.FormatInt(loginoutEvent.PUid, 10), properties)
}

func PushLoginOutEventData(LoginOutEventData LoginOutEventData) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(LoginOutEventData)
	if err != nil {
		return 0, err
	}

	return SendEvent("LoginOutLog", "-1", "-1", "", properties)
}

func PushKidLoginOutEventData(LoginOutEventData LoginOutEventData) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(LoginOutEventData)
	if err != nil {
		return 0, err
	}

	return SendEvent("LoginOutKidLog", "-1", "-1", "", properties)
}

func PushGuessLoginOutEventData(LoginOutEventData LoginOutEventData) (int, error) {
	properties, err := reflectUtil.BsonObjectToMap(LoginOutEventData)
	if err != nil {
		return 0, err
	}

	return SendEvent("LoginOutGuessLog", "-1", "-1", "", properties)
}
