package userref

import (
	"github.com/vhaoran/vchat/lib/ypg"
)

//代理商家
type BrokerInfo struct {
	ypg.BaseModel
	ID int64 `json:"id"`
	//主体名称
	Name string `json:"name" gorm:"name:name;type:text;null;"`
	//主体简介
	Brief string `json:"brief" gorm:"name:;type:text;null;"`
	//头像路径
	Icon string `json:"icon" gorm:"name:icon;type:text;null;"`

	//拥有者id,来自于userInfo
	OwnerID       int64  `json:"owner_id" gorm:"name:owner_id;null;"`
	OwnerUserCode string `json:"owner_user_code" gorm:"name:owner_user_code;type:varchar(100);null;"`
	OwnerNick     string `json:"owner_nick" gorm:"name:owner_nick;type:text;null;"`
	OwnerIcon     string `json:"owner_icon" gorm:"name:owner_icon;type:text;null;"`

	//提现帐号，只能来自微信或支付宝
	//0 支付宝 1：微信
	//只能是两者之一
	AccountType int    `json:"account_type" gorm:"name:account_type;null;"`
	AccountCode string `json:"account_code" gorm:"name:account_code;type:text;null;"`

	//使用平台商城 0:否 1：是
	EnableMall int `json:"enable_mall"`
	//使用平台大师 0:否 1：是
	EnableMaster int `json:"enable_master"`

	//使用平台悬赏貼
	EnablePrize int `json:"enable_prize"`
	//使用平台闪断貼
	EnableVie int `json:"enable_vie"`

	ServiceCode string `json:"service_code" gorm:"name:service_code;type:varchar(100);null;"`

	//启用状态（0：停用 1：启用,默认启用）
	Enabled int `json:"enabled"`
}

func (BrokerInfo) TableName() string {
	return "broker_info"
}
