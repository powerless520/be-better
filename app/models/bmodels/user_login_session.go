package bmodels

type UserLoginSession struct {
	SessionId       string   `json:"session_id"`
	LoginTime       int64    `json:"login_time"`
	PI  			string   `json:"pi"`
	PUID  		    int64    `json:"puid"`
	DI              string   `json:"di"`
}

type AppDataEvent struct {
	AppId         int    `bson:"appid" json:"appid"`
	PUid          int64  `bson:"puid" json:"puid"`
	Bt            int    `bson:"bt" json:"bt"`
	Timestamp     int64  `bson:"timestamp" json:"timestamp"` //发送上线时间
	Attempts      int    `bson:"attempts" json:"attempts"`
}