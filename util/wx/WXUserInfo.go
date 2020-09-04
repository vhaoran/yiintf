package wx


type WxUserInfoRef struct {
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Nickname   string   `json:"nickname"`
	Openid     string   `json:"openid"`
	Privilege  []string `json:"privilege"`
	Province   string   `json:"province"`
	Sex        int      `json:"sex"`
	UnionID    string   `json:"unionid"`
}

