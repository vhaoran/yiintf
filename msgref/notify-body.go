package msgref

type NotifyBody struct {
	//操作代码
	Code string `json:"code"`
	//操作的中文说明
	Comment string `json:"comment"`
	//用于表示具体业务操作的数据，可以为空串
	//id/remainder/etc..
	Data string `json:"data"`
}
