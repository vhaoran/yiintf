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

//----------------------------------------------------
// auth: whr  date:2020/9/1118:19--------------------------
// ####请勿擅改此功能代码####
// 用途：根据uid判断用户是否是管理员
//---------------------------------------------
const (
	IsAdmin_H_PATH = "/IsAdmin"
)

//获取用户所有好友
type (
	IsAdminService interface {
		Exec(in *IsAdminIn) (*IsAdminOut, error)
	}

	//input data
	IsAdminIn struct {
		UID int64 `json:"is_admin"`
	}

	//output data
	IsAdminOut struct {
		bool
	}

	// handler implements
	IsAdminH struct {
		base ykit.RootTran
	}
)

func (r *IsAdminH) MakeLocalEndpoint(svc IsAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*IsAdminIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *IsAdminH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(IsAdminIn), ctx, req)
}

//个人实现,参数不能修改
func (r *IsAdminH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *IsAdminOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *IsAdminH) HandlerLocal(service IsAdminService,
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
func (r *IsAdminH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		IsAdmin_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *IsAdminH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		IsAdmin_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_IsAdmin sync.Once
var local_IsAdmin_EP endpoint.Endpoint

func (r *IsAdminH) Call(in *IsAdminIn) (*IsAdminOut, error) {
	once_IsAdmin.Do(func() {
		local_IsAdmin_EP = new(IsAdminH).ProxySD()
	})
	//
	ep := local_IsAdmin_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*IsAdminOut), nil
}
