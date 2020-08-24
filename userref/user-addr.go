package userref

type UserAddr struct {
	//Id
	ID int64
	//用户id
	UID int64
	//收货人
	ContactPerson string
	//手机号
	Mobile string
	//地址
	AddrCommon
	//黙认地址
	IsDefault int
}
