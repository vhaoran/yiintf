package util

//last update date 2020-04-21
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
	WXLoginInner_H_PATH = "/WXLoginInner"
)

// postman
type (
	WXLoginInnerService interface {
		Exec(ctx context.Context, in *WXLoginInnerIn) (*ykit.Result, error)
	}

	//input  data
	WXLoginInnerIn struct {
		Code string `json:"code"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Body interface{} `json:"data"`
	//}

	// handler implements
	WXLoginInnerH struct {
		base ykit.RootTran
	}
)

func (r *WXLoginInnerH) MakeLocalEndpoint(svc WXLoginInnerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  WXLoginInner ###########")
		spew.Dump(ctx)

		in := request.(*WXLoginInnerIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *WXLoginInnerH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(WXLoginInnerIn), ctx, req)
}

//个人实现,参数不能修改
func (r *WXLoginInnerH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *WXLoginInnerH) HandlerLocal(service WXLoginInnerService,
	mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {

	ep := r.MakeLocalEndpoint(service)
	for _, f := range mid {
		ep = f(ep)
	}

	before := tran.ServerBefore(ykit.Jwt2ctx())
	opts := make([]tran.ServerOption, 0)
	opts = append(opts, before)
	opts = append(opts, options...)

	handler := tran.NewServer(
		ep,
		r.DecodeRequest,
		r.base.EncodeResponse,
		opts...)
	//handler = loggingMiddleware()
	return handler
}

//sd,proxy实现,用于etcd自动服务发现时的handler
func (r *WXLoginInnerH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		MSTAG,
		"POST",
		WXLoginInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *WXLoginInnerH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		MSTAG,
		"POST",
		WXLoginInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_WXLoginInner sync.Once
var local_WXLoginInner_EP endpoint.Endpoint

func (r *WXLoginInnerH) Call(in WXLoginInnerIn) (*ykit.Result, error) {
	once_WXLoginInner.Do(func() {
		local_WXLoginInner_EP = new(WXLoginInnerH).ProxySD()
	})
	//
	ep := local_WXLoginInner_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*ykit.Result), nil
}
