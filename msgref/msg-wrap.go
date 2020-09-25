package msgref

type MsgWrap struct {
	//To          int64 `json:"uid"`
	ContentType string      `json:"content_type"`
	Content     interface{} `json:"content"`
}
