package core

import (
	"be-better/core/global"
	"be-better/core/initialize"
)

func GlobalInit() {
	global.GlobalViper = initialize.Viper()
	global.GlobalLogger = global.Logger()
	global.GlobalDatabase = initialize.Gorm()

	RunServer()
}
