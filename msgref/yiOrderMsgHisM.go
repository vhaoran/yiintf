package msgref

import (
	"github.com/vhaoran/vchat/common/ytime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type YiOrderMsgHis struct {
	ID primitive.ObjectID `json:"id"   bson:"_id"`
	//是否客户: true: 客户 false：大师

	//用户id
	From int64 `json:"from"   bson:"from"`
	//用户眤称
	FromNick string `json:"from_nick"   bson:"from_nick"`
	//用户头像
	FromIcon string `json:"from_icon"   bson:"from_icon"`

	//用户id
	To int64 `json:"to"   bson:"to"`
	//用户眤称
	ToNick string `json:"to_nick"   bson:"to_nick"`
	//用户头像
	ToIcon string `json:"to_icon"   bson:"to_icon"`

	IDOfYiOrder string `json:"id_of_yi_order"`

	//消息时间
	CreateDate    ytime.Date `json:"create_date"   bson:"create_date"`
	CreateDateInt int64      `json:"create_date_int"   bson:"create_date_int"`

	//消息类型 0:文本 1：语音
	MsgType int `json:"msg_type"   bson:"msg_type"`
	//消息类型 0:文本 1：语音时表示为一个路径
	Msg string `json:"msg"   bson:"msg"`

	Ack bool `json:"ack"   bson:"ack"`
}
