package userref

import (
	"github.com/vhaoran/vchat/common/ytime"
)

//----------------------------------------------------
// auth: whr  date:2020/8/1215:27--------------------------
// ####请勿擅改此功能代码####
// 用途：
//---------------------------------------------
//大师基本信息
type MasterInfoRef struct {
	UID      int64  `json:"uid"`
	UserCode string `json:"user_code"`
	Nick     string `json:"nick" gorm:"name:nick;type:text;null;"`
	Icon     string `json:"icon" gorm:"name:icon;type:text;null;"`

	Brief string `json:"brief" gorm:"name:brief;type:text;null;"`

	//返点比例，值为30表示30%
	Rebate float64 `json:"rebate" gorm:"name:rebate;null;"`

	Rate       int64      `json:"rate"`
	SignDate   ytime.Date `json:"sign_date"`
	BestRate   int64      `json:"best_rate"`
	OrderTotal int64      `json:"order_total"`
	MidRate    int64      `json:"mid_rate"`
	BadRate    int64      `json:"bad_rate"`
	Level      int64      `json:"level"`
	Enabled    int        `json:"enabled"`
}
