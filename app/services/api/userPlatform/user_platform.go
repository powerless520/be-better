package userPlatform

import (
	"be-better/app/models"
	"be-better/utils/encryptUtil"
	"be-better/utils/netUtil"
	"be-better/utils/strUtil"
	"errors"
	"strings"
)

func NotifyTo(url string, appid int, signKey string, user *models.User, realname string, uid string) (string, error) {

	userStatus := user.Status

	if realname != "" && realname != user.RealName {
		user.Status = 2
	}

	params := map[string]interface{}{
		"appid":  appid,
		"puid":   user.PUid,
		"pi":     *user.PI,
		"status": userStatus,
		"uid":    uid,
	}

	signStr := encryptUtil.GetSignStrInterface(params, true, "&", "=")
	params["sign"] = encryptUtil.Md5Encode(signStr + "#" + signKey)

	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}

	url += strUtil.HttpBuildQuery(params)

	resp, err := netUtil.HttpGet(url, nil)
	if err != nil {
		return resp, err
	}

	if resp != "true" {
		return resp, errors.New("notify result is not excepted: " + url + ", result: " + resp)
	}

	return resp, nil
}
