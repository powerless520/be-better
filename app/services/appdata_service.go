package services

import (
	"be-better/app/models"
	"be-better/core/global"
	"context"
	"strconv"
)

const AppDAtaCacheKey = "facm:apps:data"

func InitAppsData()  {
	apps,err := GetApps()
	if err != nil{
		return
	}
	initAppDataUser(apps)
	initUserEvent(apps)
}

func initAppDataUser(apps []models.App)  {
	userMin,userMax,err := getAppDataUserNumRange()
	if err != nil{
		return
	}

	for _,app := range apps{
		appUsersKey := AppDAtaCacheKey + strconv.Itoa(app.AppId) + ":users"
		existKey,err := global.GlobalRedis.Exists(context.Background(),appUsersKey).Result()
		if err != nil{
			continue
		}
		llen,err := global.GlobalRedis.LLen(context.Background(),appUsersKey).Result()
		if err != nil{

		}


	}

}
