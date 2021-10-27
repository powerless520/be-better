package fcm

import (
	"be-better/core/global"
	"be-better/utils/encryptUtil"
	"be-better/utils/netUtil"
	"be-better/utils/strUtil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"strings"
	"time"
)

type FcmApi struct {
	AppId     string
	BizId     string
	SecretKey string
}

const IDCARD_AUTHENTICATION_CHECK_URL = "https://api.wlc.nppa.gov.cn/idcard/authentication/check"
const IDCARD_AUTHENTICATION_QUERY_URL = "http://api2.wlc.nppa.gov.cn/idcard/authentication/query"
const LOGINOUT_URL = "http://api2.wlc.nppa.gov.cn/behavior/collection/loginout"

func GetFcm(appId string, bizId string, secretKey string) (FcmApi, error) {
	return FcmApi{
		AppId:     appId,
		BizId:     bizId,
		SecretKey: secretKey,
	}, nil
}

func (this FcmApi) IdcardAuthenticationCheck(url string, ai string, name string, idNum string) (*IdcardAuthenticationResponse, error) {
	if url == "" {
		url = IDCARD_AUTHENTICATION_CHECK_URL
	}

	err := this.RetryLimit(url, "300-s")
	if err != nil {
		return nil, errors.New("request limit error")
	}

	params := map[string]interface{}{
		"ai": ai,
	}
	var body string
	for i := 0; i < 3; i++ {
		body, err = this.get(url, params)
		if err != nil {
			time.Sleep(time.Duration(10) * time.Microsecond)
			continue
		}
		break
	}
	if err != nil {
		return nil, errors.New("IdcardAuthencationQuery request error,please retry later")
	}

	var idcardAuthenticationResponse IdcardAuthenticationResponse
	err = json.Unmarshal([]byte(body), &idcardAuthenticationResponse)
	if err != nil {
		return nil, err
	}
	return &idcardAuthenticationResponse, nil
}

func (this FcmApi) IdcardAuthenticationQuery(url string, ai string) (*IdcardAuthenticationResponse, error) {

	if url == "" {
		url = IDCARD_AUTHENTICATION_QUERY_URL
	}

	err := this.RetryLimit(url, "300-s")
	if err != nil {
		return nil, errors.New("request limit error")
	}

	params := map[string]interface{}{
		"ai": ai,
	}

	var body string
	for i := 0; i < 3; i++ {
		body, err = this.get(url, params)
		if err != nil {
			time.Sleep(time.Duration(10) * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		return nil, errors.New("IdcardAuthenticationQuery request error, please retry later")
	}

	var idcardAuthenticationResponse IdcardAuthenticationResponse
	err = json.Unmarshal([]byte(body), &idcardAuthenticationResponse)
	if err != nil {
		return nil, err
	}

	return &idcardAuthenticationResponse, nil
}

func (this FcmApi) LoginOutEvent(url string, msg []LoginOutEventMessage) (*LoginOutEventResponse, error) {

	if url == "" {
		url = LOGINOUT_URL
	}

	err := this.RetryLimit(url, "10-s")
	if err != nil {
		return nil, errors.New("request limit error")
	}

	params := map[string]interface{}{
		"collections": msg,
	}

	var body string
	for i := 0; i < 3; i++ {
		body, err = this.postEncrypt(url, params)
		if err != nil {
			time.Sleep(time.Duration(10) * time.Microsecond)
			continue
		}
		break
	}

	if err != nil {
		return nil, errors.New("LoginOutEvent request error")
	}

	var loginOutEventResponse LoginOutEventResponse
	err = json.Unmarshal([]byte(body), &loginOutEventResponse)
	if err != nil {
		return nil, err
	}
	return &loginOutEventResponse, nil
}

func (this FcmApi) getLimit(key string, rateStr string) (*limiter.Context, error) {

	conmonMsg := key + "[fcm:]" + this.AppId + "@@" + this.BizId + "]"

	store, err := redis.NewStore(global.GlobalRedis)
	if err != nil {
		global.GlobalLogger.Error("request "+conmonMsg+" limit error: ", err)
		return nil, err
	}

	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		return nil, err
	}

	limitContext, err := store.Get(context.Background(), ":fcm:"+this.AppId+"@@"+this.BizId+":"+key, rate)
	if err != nil {
		return nil, err
	}

	if limitContext.Reached {
		return &limitContext, errors.New("request limit reache")
	}

	return &limitContext, nil
}

func (this FcmApi) RetryLimit(key string, rateStr string) (err error) {
	index := 0
	for index <= 5 {
		_, err = this.getLimit(key, rateStr)
		if err != nil {
			return nil
		}
		index += 1
		time.Sleep(time.Second)
	}
	return err
}

func (this FcmApi) postEncrypt(url string, params map[string]interface{}) (string, error) {
	dataBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	dataStr := string(dataBytes)
	encryptData, err := encryptUtil.AesGcmEncrypt(dataStr, this.SecretKey)
	if err != nil {
		return "", err
	}

	encryptMap := make(map[string]interface{})
	encryptMap["data"] = encryptData

	encryptBytes, err := json.Marshal(encryptMap)
	if err != nil {
		return "", err
	}

	postData := string(encryptBytes)

	headers := make(map[string]interface{})

	sysParams := map[string]interface{}{
		"appId":      this.AppId,
		"bizIz":      this.BizId,
		"timestamps": time.Now().UnixNano() / 1e6,
	}

	for k, v := range sysParams {
		headers[k] = v
	}

	signData := map[string]string{}
	for k, v := range sysParams {
		signData[k] = fmt.Sprint(v)
	}

	signStr := encryptUtil.GetSignStr(signData, true, "", "") + postData

	headers["sign"] = encryptUtil.SHA256Encode(this.SecretKey + signStr)

	return netUtil.HttpPostWithHeaderParamsJson(url, headers, postData)
}

func (this FcmApi) get(url string, params map[string]interface{}) (string, error) {

	headers := make(map[string]interface{})

	sysParams := map[string]interface{}{
		"appId":     this.AppId,
		"bizId":     this.BizId,
		"timestamp": time.Now().UnixNano() / 1e6,
	}

	for k, v := range sysParams {
		headers[k] = v
	}

	signData := map[string]string{}
	for k, v := range sysParams {
		signData[k] = fmt.Sprint(v)
	}

	for k, v := range params {
		signData[k] = fmt.Sprint(v)
	}

	signStr := encryptUtil.GetSignStr(signData, true, "", "")
	headers["sign"] = encryptUtil.SHA256Encode(this.SecretKey + signStr)

	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += strUtil.HttpBuildQuery(params)

	return netUtil.HttpGet(url, headers)
}
