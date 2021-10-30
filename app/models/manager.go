package models

import "time"

type Manager struct {
	Id              int       `json:"id" form:"id" gorm:"column:id;primary_key;"`
	Email           string    `json:"email" form:"email" gorm:"column:email;"`
	Name            string    `json:"name" form:"name" gorm:"column:name;"`
	Status          int       `json:"status" form:"status" gorm:"column:status;"`
	Token           string    `json:"token" form:"token" gorm:"column:token;"`
	TokenExpireTime time.Time `json:"token_expire_time" form:"token_expire_time" gorm:"column:token_expire_time;"`
	LastLoginIp     string    `json:"last_login_ip" form:"last_login_ip" gorm:"column:last_login_ip;"`
	CreatedAt       time.Time `json:"created_at" form:"created_at" gorm:"column:created_at;"`
	UpdatedAt       time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at;"`
}

func (Manager) TableName() string {
	return "facm.manager"
}
