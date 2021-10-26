package global

import (
	"be-better/config"
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
)
