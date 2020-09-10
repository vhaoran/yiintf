package userref

type UserAddr struct {
	//Id
	ID int64 `json:"id"`
	//用户id
	UID int64 `json:"uid" gorm:"name:uid;null;"`
	//收货人
	ContactPerson string `json:"contact_person"`
	//手机号
	Mobile string `json:"mobile"`
	//地址
	AddrCommon
	//黙认地址
	IsDefault int `json:"is_default"`
}

func (UserAddr) TableName() string {
	return "user_addr"
}
