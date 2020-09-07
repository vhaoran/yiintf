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

	"github.com/vhaoran/yiintf/userref"
)

const (
	InnerProductInfoGet_H_PATH = "/InnerProductInfoGet"
)

//获取用户所有好友
type (
	InnerProductInfoGetService interface {
		Exec(in *InnerProductInfoGetIn) (*InnerProductInfoGetOut, error)
	}

	//input data
	InnerProductInfoGetIn struct {
		IDOfES string `json:"id_of_es"`
	}

	//output data
	InnerProductInfoGetOut struct {
		*userref.Product
	}

	// handler implements
	InnerProductInfoGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerProductInfoGetH) MakeLocalEndpoint(svc InnerProductInfoGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerProductInfoGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerProductInfoGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerProductInfoGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerProductInfoGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerProductInfoGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerProductInfoGetH) HandlerLocal(service InnerProductInfoGetService,
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
func (r *InnerProductInfoGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerProductInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerProductInfoGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerProductInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerProductInfoGet sync.Once
var local_InnerProductInfoGet_EP endpoint.Endpoint

func (r *InnerProductInfoGetH) Call(in *InnerProductInfoGetIn) (*InnerProductInfoGetOut, error) {
	once_InnerProductInfoGet.Do(func() {
		local_InnerProductInfoGet_EP = new(InnerProductInfoGetH).ProxySD()
	})
	//
	ep := local_InnerProductInfoGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerProductInfoGetOut), nil
}
