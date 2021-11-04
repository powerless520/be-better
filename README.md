be-better
├── app
│   ├── http
│   │   ├── controllers
│   │   │   ├── admin
│   │   │   │   └── app.go
│   │   │   └── v1
│   │   │       ├── baseCtrl.go
│   │   │       ├── event.go
│   │   │       └── monitor.go
│   │   └── middleware
│   ├── jobs
│   ├── models
│   │   ├── app.go
│   │   ├── bmodels
│   │   │   └── user_login_session.go
│   │   ├── manager.go
│   │   ├── response
│   │   │   ├── idcard.go
│   │   │   └── pagination.go
│   │   ├── user.go
│   │   └── workerNode.go
│   └── services
│       ├── api
│       │   ├── dana
│       │   │   ├── dana.go
│       │   │   └── model.go
│       │   ├── fcm
│       │   │   ├── fcm.go
│       │   │   └── model.go
│       │   ├── manager_service.go
│       │   └── userPlatform
│       │       ├── model.go
│       │       └── user_platform.go
│       ├── app_service.go
│       ├── appdata_service.go
│       ├── base_service.go
│       ├── loginout_service.go
│       └── user_service.go
├── assets
│   └── config
│       └── application.yaml
├── bin
│   └── control.sh
├── config
│   ├── Dana.go
│   ├── config.go
│   ├── kafka.go
│   ├── mysql.go
│   ├── privacyEncrypt.go
│   ├── redis.go
│   └── system.go
├── core
│   ├── core.go
│   ├── global
│   │   ├── global.go
│   │   ├── logger.go
│   │   └── model
│   │       ├── dana.go
│   │       └── worker.go
│   ├── initialize
│   │   ├── gorm.go
│   │   ├── logger.go
│   │   └── viper.go
│   ├── queue
│   │   ├── job.go
│   │   ├── kafkaConsume.go
│   │   ├── kafkaQueue.go
│   │   ├── queue.go
│   │   └── redisQueue.go
│   ├── response
│   │   └── response.go
│   ├── server.go
│   └── worker
│       └── dbWorkerIdAssigner.go
├── desc.md
├── go.mod
├── go.sum
├── main.go
├── router
│   └── regroute.go
├── test
│   ├── 1.lua
│   ├── demo1.go
│   └── order_test.go
└── utils
├── aesUtil.go
├── algorithmUtil.go
├── commonUtil.go
├── dateUtil
│   ├── dateUtil.go
│   └── dateUtil_test.go
├── dbutil
│   └── dbutil.go
├── diffUtil.go
├── directory.go
├── encryptUtil
│   ├── base32Util.go
│   ├── encryptUtil.go
│   ├── md5Util.go
│   ├── packUtil.go
│   ├── shaUtil.go
│   ├── signUtil.go
│   └── sortUtil.go
├── fileUtil.go
├── hmacUtil.go
├── idGenerator
│   ├── idGenerator.go
│   ├── snowFlake
│   │   └── snowFlake.go
│   └── worker
│       └── workIdAssigner.go
├── mapUtil.go
├── netUtil
│   ├── netUtil.go
│   └── request.go
├── netUtil.go
├── newCommonUtil.go
├── poolUtil.go
├── randUtil
│   └── randUtil.go
├── redisutil
│   └── redisutil.go
├── reflectUtil
│   └── reflectUtil.go
├── rsaUtil.go
├── shaUtil.go
├── signUtil.go
├── strUtil
│   ├── codeUtil.go
│   └── strUtil.go
├── strUtil.go
└── xml.go

40 directories, 92 files
