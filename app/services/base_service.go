package services

import "be-better/core/global"

type BaseService struct {
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
func (a BaseService) LoggerInfo(msg ...interface{}) {
	LoggerInfo(msg)
}

// LoggerFatal ...
func (a BaseService) LoggerFatal(msg ...interface{}) {
	LoggerFatal(msg)
}

// LoggerError ...
func (a BaseService) LoggerError(msg ...interface{}) {
	LoggerError(msg)
}