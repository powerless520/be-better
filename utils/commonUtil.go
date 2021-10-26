package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func UriToMap(uri string) (params map[string]string) {
	m := make(map[string]string)
	if len(uri) < 1 { // 空字符串
		return params
	}

	indexNum := strings.Index(uri, "?")

	fmt.Println(indexNum)
	urlStr := uri
	if indexNum == -1 {
		urlStr = uri
	} else {
		urlStr = uri[indexNum:]
	}

	if urlStr[0:1] == "?" { // 有没有包含？,有的话忽略。
		urlStr = urlStr[1:]
	}

	pars := strings.Split(urlStr, "&")
	for _, par := range pars {
		parkv := strings.Split(par, "=")
		if parkv[0] != "action" {
			enEscapeUrl, _ := url.QueryUnescape(parkv[1])
			m[parkv[0]] = enEscapeUrl // 等号前面是key,后面是value
		}
	}
	return m
}

func GetSignStr(data map[string]string, includeEmptyParam bool, joinSep string) string {
	sortedKeys := Ksort(data)
	values := make([]string, 0)

	for _, v := range sortedKeys {
		if includeEmptyParam || data[v] != "" {
			values = append(values, v+"="+data[v])
		}
	}

	signStr := strings.Join(values, joinSep)
	return signStr
}

func GetSignStrWithoutKey(data map[string]string, includeEmptyParam bool, joinSep string) string {
	sortedKeys := Ksort(data)
	values := make([]string, 0)

	for _, v := range sortedKeys {
		if includeEmptyParam || data[v] != "" {
			values = append(values, data[v])
		}
	}

	signStr := strings.Join(values, joinSep)
	return signStr
}

func Ksort(data map[string]string) []string {
	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func Md5Encode(msg string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(msg))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func KsortStrHttpBuildQuery(params map[string]string) (str string) {
	var key_slice []string
	for k, _ := range params {
		key_slice = append(key_slice, k)
	}
	sort.Strings(key_slice)
	for _, k := range key_slice {
		if k != "sign" {
			str += "&" + k + "=" + params[k]
		}
	}
	str = strings.TrimLeft(str, "&")
	return
}

func ReqSignGet(params map[string]string) (sign_str, sign_new string) {
	var err error
	params_length := len(params)
	val_ary := make([]string, params_length)
	key_ary := make([]string, params_length)
	var i int = 0
	for k, v := range params {
		if k == "app_callback_url" {
			v, err = url.QueryUnescape(v)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if v != "" && k != "sign" && k != "key" && k != "paid" && k != "action" && k != "resource_id" && k != "extra_currency" && k != "cash_type" && k != "callback_url" && k != "jump_url" && k != "app_name" && k != "app_user_name" && k != "product_name" && k != "user_ip" && k != "userip" && k != "xyzs_order_time" && k != "xyzs_deviceid" {
			val_ary[i] = v
			key_ary[i] = k
		}
		i++
	}

	sort.Strings(val_ary)
	sort.Strings(key_ary)
	fmt.Println(key_ary)
	sign_str = "582df15de91b3f12d8e710073e43f4f8" + strings.Join(val_ary, "")
	sign_new = Md5Encode(sign_str)
	return
}

func UriParseMap(uri string) (params map[string]string) {
	m := make(map[string]string)
	if len(uri) < 1 { // 空字符串
		return params
	}

	indexNum := strings.Index(uri, "?")

	urlStr := uri
	if indexNum == -1 {
		urlStr = uri
	} else {
		urlStr = uri[indexNum:]
	}

	if urlStr[0:1] == "?" { // 有没有包含？,有的话忽略。
		urlStr = urlStr[1:]
	}

	pars := strings.Split(urlStr, "&")
	for _, par := range pars {
		parkv := strings.Split(par, "=")
		if parkv[0] != "action" {
			enEscapeUrl, _ := url.QueryUnescape(parkv[1])
			m[parkv[0]] = enEscapeUrl // 等号前面是key,后面是value
		}
	}
	return m
}

// GetTableNameFromOrderId 根据订单号计算order表名
func GetTableNameFromOrderId(order_id string) string {
	unixtime_str_int64 := time.Now().Unix()
	var err error
	if err = VerifyValidOrderId(order_id); err == nil {
		unixtime_str := order_id[len(order_id)-10 : len(order_id)]
		unixtime_str_int64, err = strconv.ParseInt(unixtime_str, 10, 64)
		if err != nil {

			unixtime_str_int64 = time.Now().Unix()
		}
	}
	unixtime := time.Unix(unixtime_str_int64, 0)
	year, week := unixtime.ISOWeek()
	year_week := strconv.Itoa(year*100 + week)
	return "t_order_" + year_week
}

func VerifyValidOrderId(order_id string) (err error) {
	if order_id == "" {
		return errors.New("订单号为空")
	}
	if len(order_id) < 10 {
		return errors.New("订单号错误")
	}
	return nil
}

func RequestSignGet(params map[string]string) (sign string) {
	var err error
	params_length := len(params)
	val_ary := make([]string, params_length)
	var i int = 0
	for k, v := range params {
		if k == "app_callback_url" {
			v, err = url.QueryUnescape(v)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if v != "" && k != "sign" && k != "key" && k != "paid" && k != "action" && k != "resource_id" && k != "extra_currency" &&
			k != "cash_type" && k != "callback_url" && k != "jump_url" && k != "app_name" && k != "app_user_name" &&
			k != "product_name" && k != "user_ip" && k != "userip" && k != "xyzs_order_time" && k != "xyzs_deviceid" {
			val_ary[i] = v
		}
		i++
	}

	sort.Strings(val_ary)
	sign_str := strings.Join(val_ary, "")
	fmt.Println("sign_to_str :" + sign_str)
	sign = Md5Encode("582df15de91b3f12d8e710073e43f4f8" + sign_str)
	return sign
}

func RequestPostAll(ctx *gin.Context) (params map[string]string) {
	params = map[string]string{}
	for k, v := range ctx.Request.PostForm {
		params[k] = strings.Join(v, "|")
	}

	return params
}

func KsortWithUrlEncode(params map[string]string) (str string) {
	var key_slice []string
	for k, _ := range params {
		key_slice = append(key_slice, k)
	}
	sort.Strings(key_slice)
	for _, k := range key_slice {
		if k != "sign" {
			str += "&" + k + "=" + url.QueryEscape(params[k])
		}
	}
	str = strings.TrimLeft(str, "&")
	return
}

func RequestParamsGet(ctx *gin.Context) (params map[string]string, err error) {
	params = map[string]string{}
	rawQuery := ctx.Request.URL.RawQuery
	if rawQuery == "" {
		return params, errors.New("ParamsInvalid")
	}
	m, err := url.ParseQuery(rawQuery)
	if err != nil {
		return params, errors.New("ParseQueryError")
	}
	for k, v := range m {
		params[k] = v[0]
	}

	return params, nil
}
