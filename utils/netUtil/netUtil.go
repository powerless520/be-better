package netUtil

import (
	"bytes"
	"encoding/json"
	"errors"
	"facm/core/global"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HttpGet ...
func HttpGet(url string, headers map[string]interface{}) (data string, err error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", errors.New("ErrorWhenHttpGet:" + err.Error())
	}

	//增加header选项
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, fmt.Sprint(v))
		}
	}
	start := time.Now().UnixNano() / 1e6
	//处理返回结果
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ErrorWhenHttpGet:" + err.Error())
	}

	end := time.Now().UnixNano() / 1e6
	global.GVA_LOG.Debugf("request url: %s,  %d ms", url, end-start)

	return string(body), nil
}
func HttpPost(url, post_data string, headers map[string]interface{}) (data string, err error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("POST", url, strings.NewReader(post_data))

	if err != nil {
		return "", errors.New("ErrorWhenHttpPost:" + err.Error())
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//增加header选项
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, fmt.Sprint(v))
		}
	}

	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ErrorWhenHttpPost:" + err.Error())
	}

	return string(body), nil
}
func HttpPostJson(request_url, post_data string, headers map[string]interface{}) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", request_url, reader)
	if err != nil {
		return "", err
	}
	//增加header选项
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, fmt.Sprint(v))
		}
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("HttpPostJson ResponseStatusNot200:" + resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}
func HttpPostMapJson(request_url string, post_data map[string]interface{}, headers map[string]interface{}) (data string, err error) {
	bytesData, err := json.Marshal(post_data)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", request_url, reader)
	if err != nil {
		return "", err
	}
	//增加header选项
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, fmt.Sprint(v))
		}
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}
func HttpPostWithHeader(uri string, header, post_data map[string]interface{}) (data string, err error) {
	postdata_str := ""
	if len(post_data) > 0 {
		for k, v := range post_data {
			postdata_str += "&" + k + "=" + fmt.Sprint(v)
		}
		postdata_str = postdata_str[1:]
	}
	reader := bytes.NewReader([]byte(postdata_str))
	request, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(header) > 0 {
		for k, v := range header {
			request.Header.Set(k, fmt.Sprint(v))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return string(respBytes), errors.New("statusNot200:" + resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}
func HttpGetWithHeader(uri string, header, post_data map[string]interface{}) (data string, err error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	qry := request.URL.Query()
	if len(post_data) > 0 {
		for k, v := range post_data {
			qry.Add(k, fmt.Sprint(v))
		}
		request.URL.RawQuery = qry.Encode()
	}
	if len(header) > 0 {
		for k, v := range header {
			request.Header.Set(k, fmt.Sprint(v))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(string(respBytes))
	}
	return string(respBytes), nil
}

func HttpPostWithHeaderJson(url string, sign, post_data string) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("sign", sign)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("HttpPostWithHeaderJson ResponseStatusNot200:" + resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}

func HttpPostWithHeaderParamsJson(url string, header_map map[string]interface{}, post_data string) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(header_map) > 0 {
		for k, v := range header_map {
			request.Header.Set(k, fmt.Sprint(v))
		}
	}
	start := time.Now().UnixNano() / 1e6
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("HttpPostWithHeaderJson ResponseStatusNot200:" + resp.Status + "|" + string(respBytes))
	}

	end := time.Now().UnixNano() / 1e6
	global.GVA_LOG.Debugf("request url: %s,  %d ms", url, end-start)

	return string(respBytes), nil
}
