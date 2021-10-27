package dana

type IdcardAuthenticationRequestEvent struct {
	AppId         int    `bson:"appid" json:"appid"`
	OfficialAppId string `bson:"official_appid" json:"official_appid"`
	OfficialBizId string `bson:"official_bizid" json:"official_bizid"`
	Idcard        string `bson:"idcard" json:"idcard"`
	Realname      string `bson:"realname" json:"realname"`
	Uid           string `bson:"uid" json:"uid"`
	Rid           string `bson:"rid" json:"rid"`
	NotifyUrl     string `bson:"notify_url" json:"notify_url"`
	PackageName   string `bson:"package_name" json:"package_name"`
	AppVersion    string `bson:"app_version" json:"app_version"`
}

type IdcardAuthenticationEvent struct {
	AppId         int    `bson:"appid" json:"appid"`
	OfficialAppId string `bson:"official_appid" json:"official_appid"`
	OfficialBizId string `bson:"official_bizid" json:"official_bizid"`
	Idcard        string `bson:"idcard" json:"idcard"`     //注意，该指是加密后的值
	Realname      string `bson:"realname" json:"realname"` //注意，该指是加密后的值
	Uid           string `bson:"uid" json:"uid"`
	Rid           string `bson:"rid" json:"rid"`
	NotifyUrl     string `bson:"notify_url" json:"notify_url"`
	Attempts      int    `bson:"attempts" json:"attempts"`
}

type IdcardAuthenticationQueryEvent struct {
	AppId         int    `bson:"appid" json:"appid"`
	OfficialAppId string `bson:"official_appid" json:"official_appid"`
	OfficialBizId string `bson:"official_bizid" json:"official_bizid"`
	PUid          int64  `bson:"puid" json:"puid"`
	Ai            string `bson:"ai" json:"ai"`
	Idcard        string `bson:"idcard" json:"idcard"`
	Realname      string `bson:"realname" json:"realname"`
	Uid           string `bson:"uid" json:"uid"`
	NotifyUrl     string `bson:"notify_url" json:"notify_url"`
	Attempts      int    `bson:"attempts" json:"attempts"`
}

type IdcardAuthenticationNotifyEvent struct {
	AppId     int    `bson:"appid" json:"appid"`
	PUid      int64  `bson:"puid" json:"puid"`
	Pi        string `bson:"pi" json:"pi"`
	NotifyUrl string `bson:"notify_url" json:"notify_url"`
	Idcard    string `bson:"idcard" json:"idcard"`
	Realname  string `bson:"realname" json:"realname"`
	Uid       string `bson:"uid" json:"uid"`
}

type LoginOutEvent struct {
	AppId         int    `bson:"appid" json:"appid"`
	OfficialAppId string `bson:"official_appid" json:"official_appid"`
	OfficialBizId string `bson:"official_bizid" json:"official_bizid"`
	PUid          int64  `bson:"puid" json:"puid"`
	Bt            int    `bson:"bt" json:"bt"`
	Ot            int64  `bson:"ot" json:"ot"`
	Di            string `bson:"di" json:"di"`
	Su            string `bson:"su" json:"su"` 			    //会话身份标识, 外部控制会话生成的维度。
	Id            string `bson:"id" json:"id"`               //消息ID，系统自动生成，在发送到防沉迷的时候，如果需要生成新会话的时候，会用这个ID。这样做的目的是为了通过dana里可以大概知道发送的会话可能是什么。
	Timestamp     int64  `bson:"timestamp" json:"timestamp"` //接受消息的时间
	PackageName   string `bson:"package_name" json:"package_name"`
	AppVersion    string `bson:"app_version" json:"app_version"`
}

type LoginOutEventData struct {
	AppId         int    `bson:"appid" json:"appid"`
	OfficialAppId string `bson:"official_appid" json:"official_appid"`
	OfficialBizId string `bson:"official_bizid" json:"official_bizid"`
	PUid          int64  `bson:"puid" json:"puid"`

	Si          string `bson:"si" json:"si"`
	Bt          int    `bson:"bt" json:"bt"`
	Ot          int64  `bson:"ot" json:"ot"`
	Ct          int    `bson:"ct" json:"ct"`
	Di          string `bson:"di" json:"di"`
	Pi          string `bson:"pi" json:"pi"`
	Su            string `bson:"su" json:"su"` 			    //会话身份标识, 外部控制会话生成的维度。
	Timestamp   int64  `bson:"timestamp" json:"timestamp"` //接受消息的时间
	PackageName string `bson:"package_name" json:"package_name"`
	AppVersion  string `bson:"app_version" json:"app_version"`
}
