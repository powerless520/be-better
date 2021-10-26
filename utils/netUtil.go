package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HttpGet ...
func HttpGet(url string) (data string, err error) {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := c.Get(url)
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
	return string(body), nil
}
func HttpPost(url, post_data string) (data string, err error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(post_data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ErrorWhenHttpPost:" + err.Error())
	}
	return string(body), nil
}

//func GetClientIp(orderProcess *gin.Context) (client_ip string) {
//	client_ip = orderProcess.Request.Header.Get("X-Forwarded-For")
//	if client_ip == "" || strings.Contains(client_ip, "127.0.0.1") {
//		client_ip = orderProcess.Request.Header.Get("X-real-ip")
//	}
//	if client_ip == "" {
//		client_ip = "127.0.0.1"
//	} else if strings.Contains(client_ip, ",") {
//		client_ip_ary := strings.Split(client_ip, ",")
//		client_ip = client_ip_ary[0]
//	} else if strings.Contains(client_ip, "|") {
//		client_ip_ary := strings.Split(client_ip, "|")
//		client_ip = client_ip_ary[0]
//	}
//	return client_ip
//}
func HttpPostJson(request_url, post_data string) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", request_url, reader)
	if err != nil {
		return "", err
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
func HttpPostMapJson(request_url string, post_data map[string]string) (data string, err error) {
	bytesData, err := json.Marshal(post_data)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", request_url, reader)
	if err != nil {
		return "", err
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
func HttpPostWithHeader(uri string, header, post_data map[string]string) (data string, err error) {
	postdata_str := ""
	if len(post_data) > 0 {
		for k, v := range post_data {
			postdata_str += "&" + k + "=" + v
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
			request.Header.Set(k, v)
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
func HttpGetWithHeader(uri string, header, post_data map[string]string) (data string, err error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	qry := request.URL.Query()
	if len(post_data) > 0 {
		for k, v := range post_data {
			qry.Add(k, v)
		}
		request.URL.RawQuery = qry.Encode()
	}
	if len(header) > 0 {
		for k, v := range header {
			request.Header.Set(k, v)
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

func HttpPostWithHeaderParamsJson(url string, header_map map[string]string, post_data string) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(header_map) > 0 {
		for k, v := range header_map {
			request.Header.Set(k, v)
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
		return "", errors.New("HttpPostWithHeaderJson ResponseStatusNot200:" + resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}

func HttpPostWithHeaderStatus(url string, header_map map[string]string, post_data string) (data string, err error) {
	reader := bytes.NewReader([]byte(post_data))
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(header_map) > 0 {
		for k, v := range header_map {
			request.Header.Set(k, v)
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("ResponseStatusError:" + resp.Status + "|" + string(respBytes))
	}
	return string(respBytes), nil
}
