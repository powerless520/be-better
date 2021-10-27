package global

import (
	"be-better/config"
	"be-better/core/global/model"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GlobalConfig   config.Server
	GlobalDatabase *gorm.DB
	GlobalRedis    *redis.Client
	GlobalViper    *viper.Viper
	GlobalLogger   *logrus.Logger
	GlobalDANA   *model.DanaClient
	GlobalIDGENERATORS  *model.IdGenerators
)
