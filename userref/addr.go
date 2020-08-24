package userref

type AddrCommon struct {
	//省
	Province string `json:"province"   bson:"province"`
	//城市
	City string `json:"city"   bson:"city"`
	//县区
	Area string `json:"area"   bson:"area"`
	//详细地址,包括街道等
	Detail string `json:"detail"   bson:"detail"`
	//邮编
	ZipCode string `json:"zipcode"   bson:"zipcode"`
}
