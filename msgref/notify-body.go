package msgref

type NotifyBody struct {
	//操作代码
	Code string `json:"code"   bson:"code"`
	//操作的中文说明
	Comment string `json:"comment"   bson:"comment"`
	//用于表示具体业务操作的数据，可以为空串
	//id/remainder/etc..
	Body string `json:"body"   bson:"body"`
}
