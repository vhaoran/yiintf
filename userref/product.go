package userref

type (
	//商品表 存于es中
	Product struct {
		CateId    int64           `json:"cate_id,omitempty"`
		CateName  string          `json:"cate_name,omitempty"`
		Name      string          `json:"name,omitempty"`
		Remark    string          `json:"remark,omitempty"`
		KeyWord   string          `json:"key_word,omitempty"`
		ImageMain string          `json:"image_main"`
		Images    []*ProductImage `json:"images"`
		Colors    []*ProductColor `json:"colors"`
		Enabled   bool            `json:"enabled,omitempty"`
	}

	ProductImage struct {
		Path   string `json:"path"`
		SortNo int    `json:"sort_no"`
	}

	ProductColor struct {
		Code  string  `json:"code"`
		Price float64 `json:"price"`
	}
)
