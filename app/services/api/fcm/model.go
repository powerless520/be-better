package fcm

type IdcardAuthenticationResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Result struct {
			Status int    `json:"status"`
			Pi     string `json:"pi"`
		} `json:"result"`
	} `json:"data"`
}

type LoginOutEventResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Results []struct {
			Status  int    `json:"no"`
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		} `json:"results"`
	} `json:"data"`
}

type LoginOutEventMessage struct {
	No int    `json:"no" bson:"no"`
	Si string `json:"si" bson:"si"`
	Bt int    `json:"bt" bson:"bt"`
	Ot int64  `json:"ot" bson:"ot"`
	Ct int    `json:"ct" bson:"ct"`
	Di string `json:"di" bson:"di"`
	Pi string `json:"pi" bson:"pi"`
}
