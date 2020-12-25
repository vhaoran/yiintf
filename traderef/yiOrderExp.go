package traderef

import (
	"github.com/vhaoran/vchat/common/ytime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//---订单的评价--------------------
//好评 experience
const EXP_BEST = 1

//中评
const EXP_MID = 2

//差评 experience
const EXP_BAD = 3

//----experi-------------------------------
//大师订单队列，for rabbit-mq-
const Q_EXP_YI_ORDER = "queue_yi_order"

//---------大师订单评价--------------------------------
type YiOrderExp struct {
	ID primitive.ObjectID `json:"id"   bson:"_id,omitempty"`
	//
	CreateDate ytime.Date `json:"create_date"   bson:"create_date"`
	//用于排序 time.unixNano
	CreateDateInt int64 `json:"create_date_int"   bson:"create_date_int"`

	//----------------------------------------------------
	//用户id
	UID int64 `json:"uid"   bson:"uid"`
	//用户眤称
	Nick string `json:"nick"   bson:"nick"`
	//用户头像
	Icon string `json:"icon"   bson:"icon"`

	//---------master-------------------------------------------
	//大师id
	MasterID int64 `json:"master_id"   bson:"master_id"`
	//大师眤称
	MasterNick string `json:"master_nick"   bson:"master_nick"`
	//大师头像
	MasterIcon string `json:"master_icon"   bson:"master_icon"`

	//--------broker--------------------------------------------
	BrokerID   int64  `json:"broker_id"   bson:"broker_id"`
	BrokerName string `json:"broker_name"   bson:"broker_name"`
	//----------------------------------------------------
	//yi-ordoer-id
	OrderID primitive.ObjectID `json:"order_id"   bson:"order_id"`
	//评价结果，好中差评
	ExpResult int `json:"exp_result"   bson:"exp_result"`
	//评价内容
	ExpText string `json:"exp_text"   bson:"exp_text"`
}
