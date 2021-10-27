package v1

import "be-better/core/global"

type BaseCtrl struct {

}

// LoggerInfo info级别日志记录
func LoggerInfo(msg ...interface{}) {
	global.GlobalLogger.Info(msg)
}

// LoggerFatal fatal级别日志记录
func LoggerFatal(msg ...interface{}) {
	global.GlobalLogger.Fatal(msg)
}

// LoggerError error级别日志记录
func LoggerError(msg ...interface{}) {
	global.GlobalLogger.Error(msg)
}

// LoggerInfo ...
func (a BaseCtrl) LoggerInfo(msg ...interface{}) {
	LoggerInfo(msg)
}

// LoggerFatal ...
func (a BaseCtrl) LoggerFatal(msg ...interface{}) {
	LoggerFatal(msg)
}

// LoggerError ...
func (a BaseCtrl) LoggerError(msg ...interface{}) {
	LoggerError(msg)
}