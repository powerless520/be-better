package initialize

import (
	"be-better/core/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func GormMysql() *gorm.DB {
	m := global.GlobalConfig.Mysql
	dsn := GetDsn()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode))
	if err != nil {
		global.GlobalLogger.Error("Mysql启动异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.GlobalLogger.Error("Mysql连接异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	}
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}

func GetDsn() string {
	m := global.GlobalConfig.Mysql
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
}

func gormConfig(mod bool) *gorm.Config {
	var configGorm = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch global.GlobalConfig.Mysql.LogZap {
	case "silent", "Silent":
		configGorm.Logger = Default.LogMode(logger.Silent)
	case "error", "Error":
		configGorm.Logger = Default.LogMode(logger.Error)
	case "warn", "Warn":
		configGorm.Logger = Default.LogMode(logger.Warn)
	case "info", "Info":
		configGorm.Logger = Default.LogMode(logger.Info)
	case "zap", "Zap":
		configGorm.Logger = Default.LogMode(logger.Info)
	default:
		if mod {
			configGorm.Logger = Default.LogMode(logger.Info)
			break
		}
		configGorm.Logger = Default.LogMode(logger.Silent)
	}
	return configGorm
}
