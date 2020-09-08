package userref

type UserAddr struct {
	//Id
	ID int64 `json:"id"`
	//用户id
	UID int64 `json:"uid" gorm:"name:uid;null;"`
	//收货人
	ContactPerson string
	//手机号
	Mobile string
	//地址
	AddrCommon
	//黙认地址
	IsDefault int `json:"is_default,omitempty"`
}

func (UserAddr) TableName() string {
	return "user_addr"
}
