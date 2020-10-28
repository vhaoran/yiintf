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
	InnerProductRef_H_PATH = "/InnerProductRef"
)

//获取用户所有好友
type (
	InnerProductRefService interface {
		Exec(in *InnerProductRefIn) (*InnerProductRefOut, error)
	}

	//input data
	InnerProductRefIn struct {
		BrokerID  int64 `json:"broker_id"`
		ProductID int64 `json:"product_id"`
	}

	//output data
	InnerProductRefOut struct {
		Ok bool `json:"ok"`
	}

	// handler implements
	InnerProductRefH struct {
		base ykit.RootTran
	}
)

func (r *InnerProductRefH) MakeLocalEndpoint(svc InnerProductRefService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerProductRefIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerProductRefH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerProductRefIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerProductRefH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerProductRefOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerProductRefH) HandlerLocal(service InnerProductRefService,
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
func (r *InnerProductRefH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerProductRef_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerProductRefH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerProductRef_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerProductRef sync.Once
var local_InnerProductRef_EP endpoint.Endpoint

func (r *InnerProductRefH) Call(in *InnerProductRefIn) (*InnerProductRefOut, error) {
	once_InnerProductRef.Do(func() {
		local_InnerProductRef_EP = new(InnerProductRefH).ProxySD()
	})
	//
	ep := local_InnerProductRef_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerProductRefOut), nil
}
