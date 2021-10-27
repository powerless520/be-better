package v1

import (
	"be-better/core/global"
	"be-better/core/initialize"
	"be-better/core/response"
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
)

func HealthStatus(c *gin.Context) {
	_, err := checkDb()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = checkRedis()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OK(c)
}

func checkDb() (bool, error) {
	url := initialize.GetDsn()
	db, err := sql.Open("mysql", url)
	if err != nil {
		return false, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func checkRedis() (bool, error) {
	if global.GlobalRedis == nil {
		return false, errors.New("redis init error")
	}

	_, err := global.GlobalRedis.Ping(context.Background()).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}
