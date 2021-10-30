package api

import (
	"be-better/app/models"
	"be-better/core/global"
	"be-better/utils/dbutil"
)

func GetManagerInfo(email string) (*models.Manager,error) {
	return getManagerInfo("email",email)
}

func GetManagerByToken(token string) (*models.Manager,error) {
	return getManagerInfo("token",token)
}

func UpdateManagerInfo(manager *models.Manager) error {
	return global.GlobalDatabase.Save(manager).Error
}

func getManagerInfo(queryName,queryValue string) (*models.Manager,error) {
	var manager models.Manager
	err := global.GlobalDatabase.Where(queryName+" = ?",queryValue).First(&manager).Error
	if err != err{
		if dbutil.IsNoRecord(err){
			return nil,nil
		}
		return &manager,err
	}
	return &manager,nil
}