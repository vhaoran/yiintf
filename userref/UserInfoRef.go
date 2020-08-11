package userref

type UserInfoRef struct {
	ID int64 `json:"id"`
	//帐号	登录依据，建议用手机号
	UserCode string `json:"user_code" gorm:"index:user_info_multi_code_mobile;type:varchar(200);not null;unique_index;"`
	//眤称
	Nick string `json:"nick" gorm:"type:varchar(100)"`
	//头像
	Icon string `json:"icon" gorm:"type:varchar(1000)"`

	//姓名
	UserName string `json:"user_name" gorm:"type:varchar(200)"`
	//状态	//	锁定时为false
	Enabled bool `json:"enabled"`
	//姓别(0,女1田,2保密)
	Sex int32 `json:"sex"`
	//出生年
	BirthYear int32 `json:"birth_year"`
	//出生月
	BirthMonth int32 `json:"birth_month"`
	//出生日
	BirthDay int32 `json:"birth_day"`
	//国家
	Country string `json:"country" gorm:"type:varchar(100)"`
	//省
	Province string `json:"province" gorm:"type:varchar(50)"`
	//城市
	City string `json:"city" gorm:"type:varchar(50)"`
	//县区
	Area string `json:"area" gorm:"type:varchar(50)"`
	//身份证号
	IdCard string `json:"id_card" gorm:"type:varchar(20)"`
	//商城地址
	BrokerID int64 `json:"broker_id" gorm:"name:broker_id;null;index;"`
}
