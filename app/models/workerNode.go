package models

type WorkerNode struct{
	Id       			int64  `json:"id" gorm:"column:id;primary_key;auto_increment;comment:自增编号"`
	Namespace  			string  `json:"realname" gorm:"column:namespace;comment:命名空间"`
	WorkId              int64  `json:"work_id" gorm:"column:work_id;type:int;comment:工作ID"`
	HostName            string  `json:"host_name" gorm:"column:host_name;comment:hostname"`
	Port                string  `json:"port" gorm:"column:port;comment:端口"`
	LaunchDate			*string  `json:"launch_date" gorm:"column:launch_date;type:timestamp;comment:上线时间"`
	AvailableAt			*string  `json:"available_at" gorm:"column:available_at;type:timestamp;comment:最后一次心跳时间"`
	CreatedAt           *string  `json:"created_at" gorm:"column:created_at;type:timestamp;comment:创建时间"`
	UpdatedAt           *string  `json:"updated_at" gorm:"column:updated_at;type:timestamp;comment:更新时间"`
}

func (WorkerNode) TableName() string {
	return "facm.worker_node"
}