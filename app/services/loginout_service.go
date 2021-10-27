package services

import (
	"be-better/app/models"
	"be-better/app/models/bmodels"
	"be-better/app/services/api/dana"
	"be-better/app/services/api/fcm"
	"be-better/core/global"
	"be-better/utils/dateUtil"
	"be-better/utils/encryptUtil"
	"be-better/utils/redisutil"
	"be-better/utils/reflectUtil"
	"be-better/utils/strUtil"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const CacheKeyLoginOutEvent = "facm:events:loginout:"
const LoginIntervalThreshold = 60 * 30

func DolLoginOut(loginOutEvenList []dana.LoginOutEvent) (int, error) {

	if loginOutEvenList == nil || len(loginOutEvenList) < 1 {
		return 1, nil
	}

	officalMap := map[string][]dana.LoginOutEventData{}

	for _, loginOutEvent := range loginOutEvenList {
		appid := loginOutEvent.AppId
		appIdStr := strconv.Itoa(appid)

		app, err := GetAppInfo(appid)
		if err != nil {
			continue
		}

		if app == nil {
			continue
		}

		if app.Status != 1 {
			continue
		}

		puid := loginOutEvent.PUid
		puidStr := strconv.FormatInt(puid, 10)

		di := loginOutEvent.Di
		if di != "" {
			di = strings.Replace(di, "-", "", -1)
			if len(di) > 32 {
				di = encryptUtil.Md5Encode(di)
			}
		}

		var user *models.User

		if puid > 0 {
			user, err = GetUserByPUid(puid)
			if err != nil || user == nil {
				continue
			}
		}

		commonMsg := "puid: " + puidStr + ", appid: " + appIdStr + ", di: " + loginOutEvent.Di + ", " + strUtil.ToJsonString(loginOutEvent)

		now := time.Now().Unix()
		if now-loginOutEvent.Timestamp >= 180 {
			continue
		}

		officalKey := appIdStr + "@@" + app.OfficialAppId + "@@" + app.OfficialBizId
		if _, ok := officalMap[officalKey]; !ok {
			officalMap[officalKey] = []dana.LoginOutEventData{}
		}

		var sessionKey string
		if puid > 0 {
			sessionKey = "puid-" + puidStr
		} else {
			sessionKey = "di-" + di
		}

		userLoginSession, err := getUserLoginSession(sessionKey,, app.OfficialAppId, app.OfficialBizId)
		if err != nil {
			continue
		}

		ct := 0
		if user == nil || user.Status != 0 {
			ct = 2
		}

		message := dana.LoginOutEventData{
			AppId:         appid,
			OfficialAppId: app.OfficialAppId,
			OfficialBizId: app.OfficialBizId,
			Bt:            loginOutEvent.Bt,
			Ot:            loginOutEvent.Timestamp,
			Ct:            ct,
			PackageName:   loginOutEvent.PackageName,
			PUid:          loginOutEvent.PUid,
			Timestamp:     loginOutEvent.Timestamp,
			Su:            loginOutEvent.Su,
		}

		if ct == 2 {
			message.Di = di
		} else {
			message.Pi = *user.PI
		}

		// 认证中
		if ct == 2 {
			sendLoginOutGuess(message)
			continue
		}

		// 未成年用户
		if dateUtil.IsKidForIdcard(user.IdCard) {
			sendLoginOutKid(message)
			continue
		}

		if message.Di == "" && message.Pi == "" {
			continue
		}

		timestamp := loginOutEvent.Timestamp
		//Bt:游戏用户行为类型， 0：下线 1：上线
		if loginOutEvent.Bt == 1 {
			if userLoginSession != nil {
				//timestamp-userLoginSession.LoginTime < LoginIntervalThreshold {
				//	continue
				//}
				continue
			}

			// 以下两种状态的时候生成会话ID

			message.Si = loginOutEvent.Id

			userLoginSession = &bmodels.UserLoginSession{
				SessionId: message.Si,
				LoginTime: timestamp,
				PI:        message.Pi,
				PUID:      puid,
				DI:        di,
			}

			err = setUserLoginSession(sessionKey,app.OfficialAppId,app.OfficialBizId,*userLoginSession)
			if err != nil{
				continue
			}
		}else {
			if userLoginSession == nil{
				continue
			}

			message.Si = userLoginSession.SessionId
			// 使用完删除redis的值
			err = deleteUserLoginSession(sessionKey,app.OfficialAppId,app.OfficialBizId)
			if err != nil{
				LoggerError("DoLoginOutEvent deleteUserLoginSession error:" +commonMsg)
				continue
			}
		}
		officalMap[officalKey] = append(officalMap[officalKey],message)

	}

	for officialKey,list := range officalMap{
		var officialList []fcm.LoginOutEventMessage
		for k,v := range list{
			var officialItem fcm.LoginOutEventMessage
			if err := reflectUtil.BsonConvert(v,&officialItem);err != nil{
				continue
			}
			officialItem.No = k + 1
			officialList = append(officialList,officialItem)
		}

		splitX := strings.Split(officialKey,"@@")
		if len(splitX) <= 2{
			continue
		}
		appidStr := splitX[0]
		appid, _ := strconv.Atoi(appidStr)

		facmService,err := GetFcmService(appid)
		if err != nil{
			continue
		}

		resp,err := facmService.LoginOutEvent("",officialList)
		if err != nil{
			continue
		}

		if resp.ErrCode != 0{
			continue
		}

		sendLoginOutLog(list)
	}

	return 1,nil
}

func sendLoginOutLog(list []dana.LoginOutEventData)  {
	if len(list) == 0{
		return
	}

	for _,v := range list{
		if _,err := dana.PushLoginOutEventData(v);err != nil{
			continue
		}
	}

	LoggerInfo("")
}

func sendLoginOutKid(data dana.LoginOutEventData)  {
	if _,err := dana.PushKidLoginOutEventData(data);err != nil{
		LoggerError("PushKidLoginOutEventData err: ",err,", data: ",data)
	}
	LoggerInfo("PushKidLoginOutEventData successful.")
}

func sendLoginOutGuess(data dana.LoginOutEventData)  {

	if _,err := dana.PushGuessLoginOutEventData(data);err != nil{
		LoggerError("PushGuessLoginOutEventData err: ",err,", data: ",data)
	}
	LoggerInfo("PushGuessLoginOutEventData successful.")
}

func getUserLoginSession(sessionUserId string,official_appid,official_bizid string) (*bmodels.UserLoginSession,error) {
	redisStr,err := global.GlobalRedis.Get(context.Background(),CacheKeyLoginOutEvent+official_appid+"@@"+official_bizid+":"+sessionUserId).Result()
	if err != nil{
		if redisutil.IsNoRecord(err){
			return nil,err
		}
	}

	session := &bmodels.UserLoginSession{}
	err = json.Unmarshal([]byte(redisStr),session)
	if err != nil{
		return nil,err
	}
	return session,nil

}

func setUserLoginSession(sessionKey string, offical_appid, offical_bizid string, userLoginSession bmodels.UserLoginSession) error {

	timer := time.Duration(168) * time.Hour // 7天
	userLoginSessionJson, err := json.Marshal(userLoginSession)
	if err != nil {
		return err
	}

	key := CacheKeyLoginOutEvent + offical_appid + "@@" + offical_bizid + ":" + sessionKey
	err = global.GlobalRedis.Set(context.Background(), key, string(userLoginSessionJson), timer).Err()

	return err
}

func deleteUserLoginSession(sessionKey string, official_appid, official_bizid string) error {
	key := CacheKeyLoginOutEvent + official_appid + "@@" + official_bizid + ":" + sessionKey
	return global.GlobalRedis.Del(context.Background(), key).Err()
}

func GeneraterFcmSessionId(appid int, puid int64) (string, error) {
	si, _ := global.GlobalIDGENERATORS.SessionIdGenerator.NextId()
	s := strconv.FormatInt(si, 10)

	result := "10"
	appidStr := rightStr(strconv.Itoa(appid),3)
	puidStr := rightStr(strconv.FormatInt(puid,10),5)
	result += appidStr + puidStr

	if len(s) > 22 {
		return "",errors.New("Generate session_id length error, " + s)
	}

	return result,nil
}

func rightStr(s string, rLen int) string {
	strLen := len(s)
	if rLen > strLen {
		rLen = strLen
	}

	start := strLen - rLen
	return s[start:]
}
