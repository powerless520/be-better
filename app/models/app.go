package models

type App struct {
	AppId             int    `json:"appid" form:"appid" gorm:"column:appid;"`
	AppKey            string `json:"appkey" form:"appkey" gorm:"column:appkey;"`
	Name              string `json:"name" from:"name" gorm:"column:name;"`
	OfficialAppId     string `json:"official_appid" form:"official_appid" gorm:"column:official_appid;"`
	OfficialBizId     string `json:"official_bizid" form:"official_bizid" gorm:"column:official_bizid;"`
	OfficialSecretKey string `json:"official_secret_key" form:"official_secret_key" gorm:"column:official_secret_key;"`
	Remark            string `json:"remark" from:"remark" gorm:"column:remark;"`
	CreatedAt         string `json:"created_at" form:"created_at" gorm:"column:created_at;comment:创建时间"`
	Status            int    `json:"status" form:"status" gorm:"column:status;comment:状态"`
	ModeType          int    `json:"mode_type" form:"mode_type" gorm:"column:mode_type;comment:类型"`
}

func (App) TableName() string {
	return "server.app"
}