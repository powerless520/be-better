package services

import (
	"be-better/app/models"
	"be-better/utils/redisutil"
	"strconv"
)

const CacheKeyUsers = "facm:users:"

func GetUserByPUid(puid int64) (*models.User,error) {
	puidStr := strconv.FormatInt(puid,10)
	// 先从缓存中获取是否存在该用户
	user,err := getUserCacheByColumn("puid",puidStr)
	if err != nil{
		if !redisutil.IsNoRecord(err){
			return user,err
		}
	}

	// 当缓存中不存在用户的时候，通过DB查询用户，如果查询到就更新到缓存中
	if user == nil{
		user,err = getDbUserByColumn("puid",puidStr)
		if err != nil{
			return user,err
		}

		if user != nil{
			err = setUserCacheByColumn(user,"puid",puidStr)
			if err != nil{
				return nil,err
			}
		}
	}

	return user,nil
}

func GetUserByIdCard(idcard string) (*models.User,error) {
	// 先从缓存中获取是否存在该用户
	user,err := getUserCacheByColumn("idcard",idcard)
	if err != nil{
		if !redisutil.IsNoRecord(err){
			return nil,err
		}
	}

	// 当缓存中不存在用户的时候，通过db查询用户，查询到更新到缓存中
	if user == nil{
		user,err = getDbUserByColumn("idcard",idcard)
		if err != nil{
			return user,err
		}

		if user != nil{
			err = setUserCacheByColumn(user,"idcard",idcard)
			if err != nil{
				return nil,err
			}
		}
	}
	return user,nil
}