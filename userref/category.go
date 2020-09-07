package userref

//商品分类
//词表存于es中
type Category struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	SortNo string `json:"sort_no"`
}

func (Category) TableName() string {
	return "category"
}
