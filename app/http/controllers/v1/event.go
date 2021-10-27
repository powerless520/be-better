package v1

import (
	"be-better/app/services"
	"be-better/core/response"
	"be-better/utils/netUtil"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type EventCtrl struct {
	BaseCtrl
}

var Event = EventCtrl{}

func (t EventCtrl) LoginOut(c *gin.Context)  {

	params := netUtil.RequestParamsPost(c)
	puidStr := params["puid"]
	btStr := params["bt"]
	otSt := params["ot"]
	di := params["di"]
	du := params["su"]
	if btStr == "" || otSt == "" || params["appid"] == "" || di == ""{
		response.FailWithMessage("appid, puid, bt, ot, di参数不能为空" , c)
		return
	}

	appid, err :=  strconv.Atoi(params["appid"])
	if err != nil{
		response.FailWithMessage("appid invalid" , c)
		return
	}

	puid := int64(0)
	if puidStr != ""{
		puid, err = strconv.ParseInt(puidStr, 10, 64)
		if err != nil{
			response.FailWithMessage("puid格式不正确" , c)
			return
		}
	}

	bt, err := strconv.Atoi(btStr)
	if err != nil{
		response.FailWithMessage("bt格式不正确" , c)
		return
	}

	ot, err := strconv.ParseInt(otSt, 10, 64)
	if err != nil{
		response.FailWithMessage("ot格式不正确" , c)
		return
	}

	if du != "" &&  (strings.Contains(du, "@") || strings.Contains(du, ":")){
		response.FailWithMessage("du不能包含字符@,:" , c)
		return
	}

	app, err := services.GetAppInfo(appid)
	if err != nil {
		response.FailWithMessage("app获取错误，请联系管理员" , c)
		return
	}
	if app == nil{
		response.FailWithMessage("app未配置，请联系管理员添加" , c)
		return
	}

	id, err :=  services.GenerateFcmSessionId(appid, puid)
	if err != nil{
		response.FailWithMessage("生成ID错误，请稍后再试" , c)
		return
	}

}