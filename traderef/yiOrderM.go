package traderef

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 0:待付款 1：已付款 3 已处理  4 已退款
const YiOrder_NoPay_0 = 0
const YiOrder_paid_1 = 1
const YiOrder_ok_3 = 3
const YiOrder_back_4 = 4

type YiOrder struct {
	//id
	ID primitive.ObjectID `json:"id"   bson:"_id"`
	//
	CreateDate time.Time `json:"create_date"   bson:"create_date"`
	//用于排序 time.unixNano
	CreateDateInt int64 `json:"create_date_int"   bson:"create_date_int"`

	//用户id
	UID int64 `json:"uid"   bson:"uid"`
	//用户代码
	UserCodeRef string `json:"user_code"   bson:"user_code"`
	//用户眤称
	NickRef string `json:"nick_ref"   bson:"nick_ref"`
	//用户头像
	IconRef string `json:"icon_ref"   bson:"icon_ref"`

	//大师id
	MasterID int64 `json:"master_id"   bson:"master_id"`
	//大师代码
	MasterUserCodeRef string `json:"master_user_code_ref"   bson:"master_user_code_ref"`
	//大师眤称
	MasterNickRef string `json:"master_nick_ref"   bson:"master_nick_ref"`
	//大师头像
	MasterIconRef string `json:"master_icon_ref"   bson:"master_icon_ref"`
	//订单类型	 0：四柱 1：六爻 3:合婚   20 其它
	OrderType int `json:"order_type"   bson:"order_type"`
	//订单内容
	Content interface{} `json:"content"   bson:"content"`
	//说明
	Comment string `json:"comment"   bson:"comment"`

	//金额
	Amt float64 `json:"amt"   bson:"amt"`
	//支付类型	0:积分付款 1：支付宝 2：微信
	PayType int `json:"pay_type"   bson:"pay_type"`
	//第三方付款单号
	TradeNo string `json:"trade_no"   bson:"trade_no"`
	// 订单状态	0:待付款 1：已付款 3 已处理  4 已退款
	Stat int `json:"stat"   bson:"stat"`
}

//六爻
type YiOrderLiuYao struct {
	//true:男性 false:女性
	IsMale bool `json:"is_male"   bson:"is_male"`

	//授卦时间
	Year   int `json:"year"   bson:"year"`
	Month  int `json:"month"   bson:"month"`
	Day    int `json:"day"   bson:"day"`
	Hour   int `json:"hour"   bson:"hour"`
	Minute int `json:"minute"   bson:"minute"`
	//
	YaoCode string `json:"yao_code"   bson:"yao_code"`
}

//四柱
type YiOrderSiZhu struct {
	//true 阳历 false:阴历
	IsSolar bool `json:"is_solar"   bson:"is_solar"`
	//姓名
	Name string `json:"name"   bson:"name"`
	//true: 男 false: 女
	IsMale bool `json:"is_male"   bson:"is_male"`

	Year   int `json:"year"   bson:"year"`
	Month  int `json:"month"   bson:"month"`
	Day    int `json:"day"   bson:"day"`
	Hour   int `json:"hour"   bson:"hour"`
	Minute int `json:"minute"   bson:"minute"`
}

//合婚
type YiOrderHeHun struct {
	//---------男方姓名
	NameMale string `json:"name_male"   bson:"name_male"`
	//true 阳历 false:阴历
	IsSolarMale bool `json:"is_solar_male"   bson:"is_solar_male"`

	YearMale   int `json:"year_male"   bson:"year_male"`
	MonthMale  int `json:"month_male"   bson:"month_male"`
	DayMale    int `json:"day_male"   bson:"day_male"`
	HourMale   int `json:"hour_male"   bson:"hour_male"`
	MinuteMale int `json:"minute_male"   bson:"minute_male"`

	//---------妇方姓名
	NameFeMale string `json:"name_female"   bson:"name_female"`
	//true 阳历 false:阴历
	IsSolarFeMale bool `json:"is_solar_female"   bson:"is_solar_female"`

	YearFeMale   int `json:"year_female"   bson:"year_female"`
	MonthFeMale  int `json:"month_female"   bson:"month_female"`
	DayFeMale    int `json:"day_female"   bson:"day_female"`
	HourFeMale   int `json:"hour_female"   bson:"hour_female"`
	MinuteFeMale int `json:"minute_female"   bson:"minute_female"`
}
