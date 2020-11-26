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
	InnerIsBrokerAdminOfBID_H_PATH = "/InnerIsBrokerAdminOfBID"
)

//获取用户所有好友
type (
	InnerIsBrokerAdminOfBIDService interface {
		Exec(in *InnerIsBrokerAdminOfBIDIn) (*InnerIsBrokerAdminOfBIDOut, error)
	}

	//input data
	InnerIsBrokerAdminOfBIDIn struct {
		UID      int64 `json:"uid"`
		BrokerID int64 `json:"broker_id"`
	}

	//output data
	InnerIsBrokerAdminOfBIDOut struct {
		Ok     bool   `json:"ok"`
		ErrStr string `json:"err_str"`
	}

	// handler implements
	InnerIsBrokerAdminOfBIDH struct {
		base ykit.RootTran
	}
)

func (r *InnerIsBrokerAdminOfBIDH) MakeLocalEndpoint(svc InnerIsBrokerAdminOfBIDService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerIsBrokerAdminOfBIDIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerIsBrokerAdminOfBIDH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerIsBrokerAdminOfBIDIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerIsBrokerAdminOfBIDH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerIsBrokerAdminOfBIDOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerIsBrokerAdminOfBIDH) HandlerLocal(service InnerIsBrokerAdminOfBIDService,
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
func (r *InnerIsBrokerAdminOfBIDH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerIsBrokerAdminOfBID_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerIsBrokerAdminOfBIDH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerIsBrokerAdminOfBID_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerIsBrokerAdminOfBID sync.Once
var local_InnerIsBrokerAdminOfBID_EP endpoint.Endpoint

func (r *InnerIsBrokerAdminOfBIDH) Call(in *InnerIsBrokerAdminOfBIDIn) (*InnerIsBrokerAdminOfBIDOut, error) {
	once_InnerIsBrokerAdminOfBID.Do(func() {
		local_InnerIsBrokerAdminOfBID_EP = new(InnerIsBrokerAdminOfBIDH).ProxySD()
	})
	//
	ep := local_InnerIsBrokerAdminOfBID_EP
	//
	result, _ := ep(context.Background(), in)

	return result.(*InnerIsBrokerAdminOfBIDOut), nil
}
