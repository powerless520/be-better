package services

import (
	"be-better/app/models"
	"be-better/app/services/api/fcm"
	"be-better/core/global"
	"be-better/utils/dbutil"
	"be-better/utils/redisutil"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

const AppCacheKey = "facem:apps:appid:"
const AppListCacheKey = "facm:apps:list:"

func GetAppInfo(appid int) (*models.App, error) {
	app, err := getAppInfoCache(appid)
	if err != nil {
		if !redisutil.IsNoRecord(err) {
			return app, err
		}
	}

	//当缓存中不存在用户的时候
	if app == nil {
		app, err = getAppInfoDB(appid)
		if err != nil {
			return app, err
		}

		if app != nil {
			err = setAppInfoCache(app)
			if err != nil {
				return nil, err
			}
		}
	}

	return app, nil
}

func GetFcmService(appid int) (fcm.FcmApi, error) {
	app, err := GetAppInfo(appid)
	if err != nil || app == nil {
		return fcm.FcmApi{}, err
	}

	return fcm.FcmApi{
		AppId:     app.OfficialAppId,
		BizId:     app.OfficialBizId,
		SecretKey: app.OfficialSecretKey,
	}, nil

}

func GetApps() ([]models.App,error) {
	key := AppListCacheKey + "2"
	redisStr,err := global.GlobalRedis.Get(context.Background(),key).Result()
	if err != nil{
		if !redisutil.IsNoRecord(err){
			return nil,err
		}
	}

	if redisStr == ""{
		var apps []models.App
		err = global.GlobalDatabase.Where("mode_type = 2").Where("status = 1").Find(&apps).Error
		if err != nil{
			if dbutil.IsNoRecord(err){
				return nil,err
			}
		}
		data,err := json.Marshal(apps)

		if err != nil{
			return nil,err
		}
		err = global.GlobalRedis.Set(context.Background(),key,string(data),24*time.Hour).Err()
		return apps,nil
	}

	var apps []models.App
	err = json.Unmarshal([]byte(redisStr),&apps)

	if err != nil{
		return nil,err
	}

	return apps,nil
}

func getAppInfoDB(appid int) (*models.App, error) {
	var app models.App
	err := global.GlobalDatabase.Where("appid = ?", appid).First(&app).Error
	if err != nil {
		if dbutil.IsNoRecord(err) {
			return nil, nil
		}
		return &app, err
	}
	return &app, nil
}

func getAppInfoCache(appid int) (*models.App, error) {
	redisStr, err := global.GlobalRedis.Get(context.Background(), AppCacheKey+strconv.Itoa(appid)).Result()
	if err != nil || redisStr == "" {
		return nil, err
	}

	app := &models.App{}
	err = json.Unmarshal([]byte(redisStr), app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func setAppInfoCache(app *models.App) error {

	data, err := json.Marshal(app)
	if err != nil {
		return err
	}

	err = global.GlobalRedis.Set(context.Background(), AppCacheKey+strconv.Itoa(app.AppId), string(data), 0).Err()

	return err
}

func DelAppInfoCache(app *models.App) {
	delAppInfoCache(app)

	delAppListCache()
}

func delAppInfoCache(app *models.App) {
	if err := global.GlobalRedis.Del(context.Background(), AppCacheKey+strconv.Itoa(app.AppId)).Err(); err != nil {
		global.GlobalLogger.Errorf("Delete app info cache,appid: %d, err: %v", app.AppId, err)
	}
}

func delAppListCache() {
	if err := global.GlobalRedis.Del(context.Background(), AppListCacheKey+"2").Err(); err != nil {
		global.GlobalLogger.Errorf("Delete app list cache, err: %v", err)
	}
}
