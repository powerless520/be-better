package strUtil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func HttpBuildQuery(data map[string]interface{}) string {
	sorted_keys := make([]string, 0)
	for k, _ := range data {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for i, k := range sorted_keys {
		value := fmt.Sprintf("%v", data[k])
		if i != (len(sorted_keys) - 1) {
			signStrings = signStrings + k + "=" + value + "&"
		} else {
			signStrings = signStrings + k + "=" + value
		}
	}
	return signStrings
}

func ToInterfaceMap(data map[string]string) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range data {
		result[k] = v
	}
	return result
}

func InterfaceToMap(data map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range data {
		result[k] = fmt.Sprint(v)
	}
	return result
}

func UriToMap(uri string) (params map[string]interface{}) {
	m := make(map[string]interface{})
	if len(uri) < 1 { // 空字符串
		return params
	}
	if uri[0:1] == "?" { // 有没有包含？,有的话忽略。
		uri = uri[1:]
	}

	pars := strings.Split(uri, "&")
	for _, par := range pars {
		parkv := strings.Split(par, "=")
		if parkv[0] != "action"{
			enEscapeUrl, _ := url.QueryUnescape(parkv[1])
			m[parkv[0]] = enEscapeUrl // 等号前面是key,后面是value
		}
	}
	return m
}

func ToJsonString(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil{
		return ""
	}
	return string(data)
}