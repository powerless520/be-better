package models

type User struct{
	PUid       			int64  `json:"puid" form:"puid" gorm:"column:puid;primary_key;comment:系统生成的用户ID"`
	RealName  			string  `json:"realname" form:"realname" gorm:"column:realname;comment:姓名"`
	IdCard              string  `json:"idcard" form:"idcard" gorm:"column:idcard;comment:身份证号"`
	Status              int     `json:"status" form:"status" gorm:"column:status;comment:状态, 0:已认证; 1: 认证中; 2: 认证失败"`
	PI                  *string  `json:"pi" form:"pi" gorm:"column:pi;comment:已通过实名认证用户的唯一标识"`
	CreatedAt           *string  `json:"created_at" form:"created_at" gorm:"column:created_at;type:timestamp;comment:创建时间"`
	UpdatedAt           *string  `json:"updated_at" form:"updated_at" gorm:"column:updated_at;type:timestamp;comment:更新时间"`
	CertificationAt     *string  `json:"certification_at" form:"certification_at" gorm:"column:certification_at;type:timestamp;comment:认证时间(认证成功)"`
}

func (User) TableName() string {
	return "facm.user"
}