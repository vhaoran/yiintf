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

	"github.com/vhaoran/yiintf/traderef"
)

const (
	InnerYiOrderGet_H_PATH = "/InnerYiOrderGet"
)

//获取用户所有好友
type (
	InnerYiOrderGetService interface {
		Exec(in *InnerYiOrderGetIn) (*InnerYiOrderGetOut, error)
	}

	//input data
	InnerYiOrderGetIn struct {
		IDHex string `json:"id"`
	}

	//output data
	InnerYiOrderGetOut struct {
		traderef.YiOrder
	}

	// handler implements
	InnerYiOrderGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerYiOrderGetH) MakeLocalEndpoint(svc InnerYiOrderGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerYiOrderGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerYiOrderGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerYiOrderGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerYiOrderGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerYiOrderGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerYiOrderGetH) HandlerLocal(service InnerYiOrderGetService,
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
func (r *InnerYiOrderGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerYiOrderGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerYiOrderGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerYiOrderGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerYiOrderGet sync.Once
var local_InnerYiOrderGet_EP endpoint.Endpoint

func (r *InnerYiOrderGetH) Call(in *InnerYiOrderGetIn) (*InnerYiOrderGetOut, error) {
	once_InnerYiOrderGet.Do(func() {
		local_InnerYiOrderGet_EP = new(InnerYiOrderGetH).ProxySD()
	})
	//
	ep := local_InnerYiOrderGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerYiOrderGetOut), nil
}
