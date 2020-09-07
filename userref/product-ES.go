package userref

import (
	"github.com/vhaoran/vchat/common/ytime"
)

type (
	//商品表 存于es中
	Product struct {
		IDOfES      string          `json:"id_of_es"`
		Created     ytime.Date      `json:"created"`
		LastUpdated ytime.Date      `json:"last_updated"`

		CateId      int64           `json:"cate_id"`
		CateName    string          `json:"cate_name"`
		Name        string          `json:"name"`
		Remark      string          `json:"remark"`
		KeyWord     string          `json:"key_word"`
		ImageMain   string          `json:"image_main"`
		Images      []*ProductImage `json:"images"`
		Colors      []*ProductColor `json:"colors"`
		Enabled     bool            `json:"enabled"`
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
