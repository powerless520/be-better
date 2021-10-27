package admin

import (
	"be-better/app/models"
	"be-better/core/global"
	"be-better/core/response"
	"be-better/utils/dbutil"
	"be-better/utils/netUtil"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AppList(ctx *gin.Context)  {
	params,err := netUtil.RequestParamsGet(ctx)
	if err != nil{
		global.GlobalLogger.Errorln(err.Error())
		return
	}

	var (
		query = "?"
		values = []interface{}{"1"}
		perPage = 20
		page  int64
		total int64
		totalPage int64
		rows []models.App
	)

	if pageStr,_ := params["page"];pageStr != ""{
		page,_ = strconv.ParseInt(pageStr,10,64)
		if page <= 0{
			page = 1
		}
	}else {
		page = 1
	}

	if perPageStr,_ := params["per_page"];perPageStr != ""{
		perPage,_ = strconv.Atoi(perPageStr)
	}

	if appid, _ := params["appid"]; appid != "" {
		query += " and appid = ?"
		values = append(values, appid)
	}

	if name, _ := params["name"]; name != "" {
		query += " and name like ?"
		values = append(values, "%"+name+"%")
	}

	if officialBizId, _ := params["official_bizid"]; officialBizId != "" {
		query += " and official_bizid = ?"
		values = append(values, officialBizId)
	}
	if statusStr, _ := params["status"]; statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query += " and status = ?"
		values = append(values, status)
	}
	if modeTypeStr, _ := params["mode_type"]; modeTypeStr != "" {
		modeType, _ := strconv.Atoi(modeTypeStr)
		query += " and mode_type = ?"
		values = append(values, modeType)
	}

	global.GlobalDatabase.Model(models.App{}).Where(query,values...).Count(&total)
	totalPage = total/int64(perPage) + 1
	if page > totalPage{
		page = totalPage
	}

	err = global.GlobalDatabase.Where(query,values...).Limit(perPage).Offset((int(page) - 1)*perPage).Order("created desc").Find(&rows).Error
	if err != nil && !dbutil.IsNoRecord(err){
		global.GlobalLogger.Errorf("Get app list err: %v",err)
		response.FailWithMessage("网络繁忙，请稍后再试",ctx)
	}

}
