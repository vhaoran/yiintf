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
	InnerMasterRef_H_PATH = "/InnerMasterRef"
)

//获取用户所有好友
type (
	InnerMasterRefService interface {
		Exec(in *InnerMasterRefIn) (*InnerMasterRefOut, error)
	}

	//input data
	InnerMasterRefIn struct {
		BrokerID int64 `json:"broker_id"`
		MasterID int64 `json:"master_id"`
	}

	//output data
	InnerMasterRefOut struct {
		Ok bool `json:"ok"`
	}

	// handler implements
	InnerMasterRefH struct {
		base ykit.RootTran
	}
)

func (r *InnerMasterRefH) MakeLocalEndpoint(svc InnerMasterRefService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerMasterRefIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerMasterRefH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerMasterRefIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerMasterRefH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerMasterRefOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerMasterRefH) HandlerLocal(service InnerMasterRefService,
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
func (r *InnerMasterRefH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerMasterRef_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerMasterRefH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerMasterRef_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerMasterRef sync.Once
var local_InnerMasterRef_EP endpoint.Endpoint

func (r *InnerMasterRefH) Call(in *InnerMasterRefIn) (*InnerMasterRefOut, error) {
	once_InnerMasterRef.Do(func() {
		local_InnerMasterRef_EP = new(InnerMasterRefH).ProxySD()
	})
	//
	ep := local_InnerMasterRef_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerMasterRefOut), nil
}
