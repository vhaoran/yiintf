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
)

const (
	InnerPayReq_H_PATH = "/InnerPayReq"
)

//获取用户所有好友
type (
	InnerPayReqService interface {
		Exec(in *InnerPayReqIn) (*InnerPayReqOut, error)
	}

	//input data
	InnerPayReqIn struct {
		IP      string  `json:"ip"`
		UID     int64   `json:"uid"`
		BType   string  `json:"b_type"`
		AccType int     `json:"acc_type"`
		TradeNo string  `json:"trade_no"`
		Amt     float64 `json:"amt"`
	}

	//output data
	InnerPayReqOut struct {
		URL string `json:"url"`
		M   ykit.M `json:"m"`
	}

	// handler implements
	InnerPayReqH struct {
		base ykit.RootTran
	}
)

func (r *InnerPayReqH) MakeLocalEndpoint(svc InnerPayReqService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerPayReqIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerPayReqH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerPayReqIn), ctx, req)
	//if err != nil{
	//	return dst,err
	//}
	//
	//if dst != nil{
	//	obj := dst.(*InnerPayReqIn)
	//	obj.IP = g.GetClientIP(req)
	//}
	//return dst,
}

//个人实现,参数不能修改
func (r *InnerPayReqH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerPayReqOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerPayReqH) HandlerLocal(service InnerPayReqService,
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
func (r *InnerPayReqH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerPayReq_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerPayReqH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerPayReq_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerPayReq sync.Once
var local_InnerPayReq_EP endpoint.Endpoint

func (r *InnerPayReqH) Call(in *InnerPayReqIn) (*InnerPayReqOut, error) {
	once_InnerPayReq.Do(func() {
		local_InnerPayReq_EP = new(InnerPayReqH).ProxySD()
	})
	//
	ep := local_InnerPayReq_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerPayReqOut), nil
}
