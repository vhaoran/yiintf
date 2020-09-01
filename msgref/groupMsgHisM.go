package msgref

import (
	"github.com/vhaoran/vchat/common/ytime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMsgHis struct {
	ID primitive.ObjectID `json:"id"   bson:"_id"`
	//是否客户: true: 客户 false：大师

	//发言人 用户id
	UID int64 `json:"uid"   bson:"uid"`
	//发言人  用户眤称
	Nick string `json:"nick"   bson:"nick"`
	//发言人  用户头像
	Icon string `json:"icon"   bson:"icon"`

	//聊天室消息（直播消息）
	//可以认为每个大师都有聊天室，故gid-master-id
	GID   int64    `json:"gid"   bson:"gid"`
	GName string `json:"g_name"   bson:"g_name"`
	GIcon string `json:"g_icon"   bson:"g_icon"`

	//消息时间
	CreateDate    ytime.Date `json:"create_date"   bson:"create_date"`
	CreateDateInt int64      `json:"create_date_int"   bson:"create_date_int"`

	//消息类型 0:文本 1：语音
	ContentType int `json:"content_type"   bson:"content_type"`
	//消息类型 0:文本 1：语音时表示为一个路径
	Content string `json:"content"   bson:"content"`
}
