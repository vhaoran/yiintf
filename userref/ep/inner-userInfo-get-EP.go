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
	InnerUserInfoGet_H_PATH = "/InnerUserInfoGet"
)

//获取用户所有好友
type (
	InnerUserInfoGetService interface {
		Exec(in *InnerUserInfoGetIn) (*InnerUserInfoGetOut, error)
	}

	//input data
	InnerUserInfoGetIn struct {
		UID    int64  `json:"uid"`
		Mobile string `json:"mobile"`
	}

	//output data
	InnerUserInfoGetOut struct {
		userref.UserInfoRef
	}

	// handler implements
	InnerUserInfoGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerUserInfoGetH) MakeLocalEndpoint(svc InnerUserInfoGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerUserInfoGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerUserInfoGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerUserInfoGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerUserInfoGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerUserInfoGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerUserInfoGetH) HandlerLocal(service InnerUserInfoGetService,
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
func (r *InnerUserInfoGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerUserInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerUserInfoGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerUserInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerUserInfoGet sync.Once
var local_InnerUserInfoGet_EP endpoint.Endpoint

func (r *InnerUserInfoGetH) Call(in *InnerUserInfoGetIn) (*InnerUserInfoGetOut, error) {
	once_InnerUserInfoGet.Do(func() {
		local_InnerUserInfoGet_EP = new(InnerUserInfoGetH).ProxySD()
	})
	//
	ep := local_InnerUserInfoGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	if result != nil {
		return result.(*InnerUserInfoGetOut), nil
	}

	return nil, nil
}
