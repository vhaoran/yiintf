package userref

//商品分类
//词表存于es中
type Category struct {
	ID     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Icon   string `json:"icon,omitempty"`
	SortNo string `json:"sort_no,omitempty"`
}
