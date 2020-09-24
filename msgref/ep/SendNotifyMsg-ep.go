package ep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"

	"github.com/vhaoran/yiintf/msgref"
)

//----------------------------------------------------
// auth: whr  date:2020/9/2411:33--------------------------
// ####请勿擅改此功能代码####
// 用途：发送内部的通知消息到前台
//--------------------------------------------- 

const (
	InnerSendNotifyMsg_H_PATH = "/InnerSendNotifyMsg"
)

//获取用户所有好友
type (
	InnerSendNotifyMsgService interface {
		Exec(in *InnerSendNotifyMsgIn) (*InnerSendNotifyMsgOut, error)
	}

	//input data
	InnerSendNotifyMsgIn struct {
		To int64 `json:"to"`
		msgref.NotifyBody
	}

	//output data
	InnerSendNotifyMsgOut struct {
		RetID  string `json:"ret_id"`
		ErrStr string `json:"err_str"`
	}

	// handler implements
	InnerSendNotifyMsgH struct {
		base ykit.RootTran
	}
)

func (r *InnerSendNotifyMsgH) MakeLocalEndpoint(svc InnerSendNotifyMsgService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerSendNotifyMsgIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerSendNotifyMsgH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerSendNotifyMsgIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerSendNotifyMsgH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerSendNotifyMsgOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerSendNotifyMsgH) HandlerLocal(service InnerSendNotifyMsgService,
	mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {

	ep := r.MakeLocalEndpoint(service)
	for _, f := range mid {
		ep = f(ep)
	}

	handler := tran.NewServer(
		ep,
		r.DecodeRequest,
		r.base.EncodeResponse,
		options...)
	//handler = loggingMiddleware()
	return handler
}

//sd,proxy实现,用于etcd自动服务发现时的handler
func (r *InnerSendNotifyMsgH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerSendNotifyMsg_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerSendNotifyMsgH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerSendNotifyMsg_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerSendNotifyMsg sync.Once
var local_InnerSendNotifyMsg_EP endpoint.Endpoint

func (r *InnerSendNotifyMsgH) Call(in *InnerSendNotifyMsgIn) (*InnerSendNotifyMsgOut, error) {
	once_InnerSendNotifyMsg.Do(func() {
		local_InnerSendNotifyMsg_EP = new(InnerSendNotifyMsgH).ProxySD()
	})
	//
	ep := local_InnerSendNotifyMsg_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerSendNotifyMsgOut), nil
}
