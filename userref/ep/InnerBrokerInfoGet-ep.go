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
	InnerBrokerInfoGet_H_PATH = "/InnerBrokerInfoGet"
)

//获取用户所有好友
type (
	InnerBrokerInfoGetService interface {
		Exec(in *InnerBrokerInfoGetIn) (*InnerBrokerInfoGetOut, error)
	}

	//input data
	InnerBrokerInfoGetIn struct {
		BrokerID int64 `json:"broker_id"`
	}

	//output data
	InnerBrokerInfoGetOut struct {
		userref.BrokerInfo
	}

	// handler implements
	InnerBrokerInfoGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerBrokerInfoGetH) MakeLocalEndpoint(svc InnerBrokerInfoGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBrokerInfoGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBrokerInfoGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBrokerInfoGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBrokerInfoGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBrokerInfoGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBrokerInfoGetH) HandlerLocal(service InnerBrokerInfoGetService,
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
func (r *InnerBrokerInfoGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBrokerInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBrokerInfoGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBrokerInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBrokerInfoGet sync.Once
var local_InnerBrokerInfoGet_EP endpoint.Endpoint

func (r *InnerBrokerInfoGetH) Call(in *InnerBrokerInfoGetIn) (*InnerBrokerInfoGetOut, error) {
	once_InnerBrokerInfoGet.Do(func() {
		local_InnerBrokerInfoGet_EP = new(InnerBrokerInfoGetH).ProxySD()
	})
	//
	ep := local_InnerBrokerInfoGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerBrokerInfoGetOut), nil
}
