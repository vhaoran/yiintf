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
	SendNotifyMsg_H_PATH = "/SendNotifyMsg"
)

//获取用户所有好友
type (
	SendNotifyMsgService interface {
		Exec(in *SendNotifyMsgIn) (*SendNotifyMsgOut, error)
	}

	//input data
	SendNotifyMsgIn struct {
		msgref.NotifyBody
	}

	//output data
	SendNotifyMsgOut struct {
		OK     bool   `json:"ok"`
		ErrStr string `json:"err_str"`
	}

	// handler implements
	SendNotifyMsgH struct {
		base ykit.RootTran
	}
)

func (r *SendNotifyMsgH) MakeLocalEndpoint(svc SendNotifyMsgService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*SendNotifyMsgIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *SendNotifyMsgH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(SendNotifyMsgIn), ctx, req)
}

//个人实现,参数不能修改
func (r *SendNotifyMsgH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *SendNotifyMsgOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *SendNotifyMsgH) HandlerLocal(service SendNotifyMsgService,
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
func (r *SendNotifyMsgH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		SendNotifyMsg_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *SendNotifyMsgH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		SendNotifyMsg_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_SendNotifyMsg sync.Once
var local_SendNotifyMsg_EP endpoint.Endpoint

func (r *SendNotifyMsgH) Call(in *SendNotifyMsgIn) (*SendNotifyMsgOut, error) {
	once_SendNotifyMsg.Do(func() {
		local_SendNotifyMsg_EP = new(SendNotifyMsgH).ProxySD()
	})
	//
	ep := local_SendNotifyMsg_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*SendNotifyMsgOut), nil
}
