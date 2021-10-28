package services

import (
	"be-better/app/models"
	response2 "be-better/app/models/response"
	"be-better/app/services/api/dana"
	"be-better/app/services/api/fcm"
	"be-better/app/services/api/userPlatform"
	"be-better/core/global"
	"be-better/core/queue"
	"be-better/utils/dateUtil"
	"be-better/utils/dbutil"
	"be-better/utils/encryptUtil"
	"be-better/utils/redisutil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const CacheKeyUsers = "facm:users:"

func GetUserByPUid(puid int64) (*models.User, error) {
	puidStr := strconv.FormatInt(puid, 10)
	// 先从缓存中获取是否存在该用户
	user, err := getUserCacheByColumn("puid", puidStr)
	if err != nil {
		if !redisutil.IsNoRecord(err) {
			return user, err
		}
	}

	// 当缓存中不存在用户的时候，通过DB查询用户，如果查询到就更新到缓存中
	if user == nil {
		user, err = getDbUserByColumn("puid", puidStr)
		if err != nil {
			return user, err
		}

		if user != nil {
			err = setUserCacheByColumn(user, "puid", puidStr)
			if err != nil {
				return nil, err
			}
		}
	}

	return user, nil
}

func GetUserByIdCard(idcard string) (*models.User, error) {
	// 先从缓存中获取是否存在该用户
	user, err := getUserCacheByColumn("idcard", idcard)
	if err != nil {
		if !redisutil.IsNoRecord(err) {
			return nil, err
		}
	}

	// 当缓存中不存在用户的时候，通过db查询用户，查询到更新到缓存中
	if user == nil {
		user, err = getDbUserByColumn("idcard", idcard)
		if err != nil {
			return user, err
		}

		if user != nil {
			err = setUserCacheByColumn(user, "idcard", idcard)
			if err != nil {
				return nil, err
			}
		}
	}
	return user, nil
}

func AddUser(user models.User) error {
	return global.GlobalDatabase.Create(&user).Error
}

func UpdateUser(user *models.User) error {
	err := global.GlobalDatabase.Save(user).Error
	if err != nil {
		return err
	}

	puid := strconv.FormatInt(user.PUid, 10)
	//清楚该用户的缓存，让下次查询的时候，能够查询到更新后的结果
	err = global.GlobalRedis.Del(context.Background(), CacheKeyUsers+"puid:"+puid, CacheKeyUsers+"idcard:"+user.IdCard).Err()
	if err != nil {
		return err
	}
	return nil
}

func UserIdcardAuthentication(idcardAuthenticationEvent dana.IdcardAuthenticationEvent, isJob bool) (*response2.AuthenticationResult, error) {
	totalStart := time.Now().UnixNano() / 1e6
	start := time.Now().UnixNano() / 1e6

	timeMsg := ""

	idcard := idcardAuthenticationEvent.Idcard
	realname := idcardAuthenticationEvent.Realname

	if idcard == "" || realname == ""{
		return nil,errors.New("参数不能为空")
	}

	notifyUrl := idcardAuthenticationEvent.NotifyUrl
	appid := idcardAuthenticationEvent.AppId

	app,err := GetAppInfo(appid)

	timeMsg += fmt.Sprintf("GetAppInfo time: %d\n",time.Now().UnixNano() / 1e6 - start)
	start = time.Now().UnixNano() / 1e6

	if err != nil || app == nil{
		return nil,errors.New("app获取错误，请联系管理员")
	}

	idcardDecrypt,err := encryptUtil.Base32Decrypt(idcard,global.GlobalConfig.PrivacyEncrypt.Key,global.GlobalConfig.PrivacyEncrypt.Iv)
	if err != nil{
		return nil,errors.New("身份证信息解密失败")
	}

	timeMsg += fmt.Sprintf("Idcard Base32Decrypt time: %d\n",time.Now().UnixNano() / 1e6 - start)
	start = time.Now().UnixNano() / 1e6

	user,err := GetUserByIdCard(idcard)
	if err != nil{
		return nil,errors.New("认证失败，请稍后重试")
	}

	timeMsg += fmt.Sprintf("GetUserByIdcard time: %d\n",time.Now().UnixNano()/1e6 - start)

	// 用户未认证过，调用服务进行认证（服务可能是异步认证）
	// 注意：
	// 1.如果服务返回限流错误，会把请求丢到队列中后续继续发送请求。

	if user != nil{

		if user.Status == 2 {

			nowStr := dateUtil.Now()
			// 当最后一次认证时间小于2小时，不继续认证

			canCheck := true
			if user.UpdatedAt !=nil{
				updateStr := dateUtil.ParseInLocationDefault(*user.UpdatedAt)
				if time.Now().Sub(updateStr).Hours() < 2{
					canCheck = false
				}
			}

			if canCheck{
				user.Status = 1
				user.RealName = realname
				user.UpdatedAt = &nowStr

				err = UpdateUser(user)
				if err != nil{
					return nil,err
				}

				err = idcardAuthenticationCheck(idcardAuthenticationEvent,*user,idcardDecrypt,realnameDecrypt,isJob)
				if err != nil{
					return nil,err
				}

			}

		}else if user.Status == 0 {

			if user.RealName != realname{
				// 返回认证失败
				result := response2.AuthenticationResult{
					Status: 2,
					PI:     user.PI,
					PUid:   strconv.FormatInt(user.PUid,10),
				}
				return &result,nil
			}

		}else if user.Status == 1 {

			if user.RealName != realname{
				return nil,errors.New("认证中，身份证和姓名不匹配")
			}

			err = idcardAuthenticationCheck(idcardAuthenticationEvent, *user, idcardDecrypt, realnameDecrypt, isJob)
			if err != nil {
				return nil, err
			}

		}

		if user.Status != 1 && notifyUrl != ""{
			err = dana.PushIdcardAuthenticationNotify()
		}

	}
}

func updateUserAuthenticationFail(attempts int,user *models.User) (bool,error) {

	if attempts < 13{
		return false,nil
	}

	authenticationResult := fcm.IdcardAuthenticationResponse{}
	authenticationResult.ErrCode = 0
	authenticationResult.Data.Result.Status = 2

	err := updateUserAuthenticationResult(user,&authenticationResult)
	if err != nil{
		return false,err
	}

	return true,nil

}

func UserIdcardAuthenticationNotify(idcardAuthenticationNotifyEvent dana.IdcardAuthenticationNotifyEvent) (string,error) {
	appid := idcardAuthenticationNotifyEvent.AppId
	app,err := GetAppInfo(appid)
	if err != nil{
		return "",errors.New("")
	}

	if app == nil{
		return "",errors.New("")
	}

	puidStr := strconv.FormatInt(idcardAuthenticationNotifyEvent.PUid,10)
	user,err := GetUserByPUid(idcardAuthenticationNotifyEvent.PUid)
	if err != nil || user == nil{
		return "",errors.New("")
	}

	if user.Status == 1{
		return "",errors.New("")
	}

	resp,err := userPlatform.NotifyTo(idcardAuthenticationNotifyEvent.NotifyUrl,appid,app.AppKey,user,idcardAuthenticationNotifyEvent.Realname,idcardAuthenticationNotifyEvent.Uid)
	if err != nil{
		return resp,err
	}

	return resp,nil

}

func GenerateFcmAI(puid int64) string {
	s := strconv.FormatInt(puid,10)
}


func PushIdcardAuthentication(idcardAuthenticationEvent *dana.IdcardAuthenticationEvent) error {
	_,err := dana.PushIdcardAuthentication(*idcardAuthenticationEvent)
	if err != nil{
		return errors.New("认证失败，稍后再试")
	}

	return nil
}

func pushIdcardAuthenticationNotify(appid int,user models.User,notifyUrl string,idcard string,realname string,uid string) error {
	notifyData := dana.IdcardAuthenticationNotifyEvent{
		AppId:     appid,
		PUid:      user.PUid,
		Pi:        *user.PI,
		NotifyUrl: notifyUrl,
		Idcard:    idcard,
		Realname:  realname,
		Uid:       uid,
	}

	_,err := dana.PushIdcardAuthenticationNotify(notifyData)
	if err != nil{
		return errors.New("认证失败，稍后再试")
	}
	return nil
}

func PushIdcardAuthenticationQuery(idcardAuthenticationQueryEvent *dana.IdcardAuthenticationQueryEvent,delaySeconds int) error {

	if err := queue.NewQueue(idcardAuthenticationQueryEvent,queue.NewJob("idcardAuthenticationQuery",*idcardAuthenticationQueryEvent)).
		SetDelay(time.Duration(delaySeconds)*time.Second).Push();err != nil{
			return err
	}
	return nil

}

func pushIdcardAuthenticationDelay(idcardAuthenticationEvent *dana.IdcardAuthenticationEvent,delaySeconds int) error{

	if err := queue.NewQueue("idcardAuthentication",queue.NewJob("idcardAuthentication",*idcardAuthenticationEvent)).
		SetDelay(time.Duration(delaySeconds) * time.Second).Push();err != nil{
			return err
	}
	return nil
}


func getDbUserByColumn(columnName string, columnValue string) (*models.User, error) {
	var user = models.User{}
	var err = global.GlobalDatabase.Where(columnName+" = ?", columnValue).First(&user).Error
	if err != nil {
		if dbutil.IsNoRecord(err) {
			return nil, nil
		}
		return &user, err
	}
	return &user, nil
}

func getUserCacheByColumn(columnName string, columnValue string) (*models.User, error) {
	userStr, err := global.GlobalRedis.Get(context.Background(), CacheKeyUsers+columnName+":"+columnValue).Result()
	if err != nil {
		return nil, err
	}

	if userStr == "" {
		return nil, nil
	}

	user := &models.User{}
	err = json.Unmarshal([]byte(userStr), user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func setUserCacheByColumn(user *models.User, columnName string, columnValue string) error {
	timer := time.Duration(30) * time.Minute
	data, err := json.Marshal(user)

	if err != nil {
		return err
	}

	err = global.GlobalRedis.Set(context.Background(), CacheKeyUsers+columnName+":"+columnValue, string(data), timer).Err()
	return err
}

func getDelayTime(attempts int) int {
	if attempts > 10{
		return 86400
	}else {
		return (1 << attempts) * 60
	}
}
