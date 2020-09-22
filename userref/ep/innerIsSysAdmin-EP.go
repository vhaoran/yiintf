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
	InnerIsSysAdmin_H_PATH = "/InnerIsSysAdmin"
)

//获取用户所有好友
type (
	InnerIsSysAdminService interface {
		Exec(in *InnerIsSysAdminIn) (*InnerIsSysAdminOut, error)
	}

	//input data
	InnerIsSysAdminIn struct {
		UID int64 `json:"is_admin"`
	}

	//output data
	InnerIsSysAdminOut struct {
		OK     bool   `json:"ok"`
		ErrStr string `json:"err_str,omitempty"`
	}

	// handler implements
	InnerIsSysAdminH struct {
		base ykit.RootTran
	}
)

func (r *InnerIsSysAdminH) MakeLocalEndpoint(svc InnerIsSysAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerIsSysAdminIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerIsSysAdminH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerIsSysAdminIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerIsSysAdminH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerIsSysAdminOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerIsSysAdminH) HandlerLocal(service InnerIsSysAdminService,
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
func (r *InnerIsSysAdminH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerIsSysAdmin_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerIsSysAdminH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerIsSysAdmin_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerIsSysAdmin sync.Once
var local_InnerIsSysAdmin_EP endpoint.Endpoint

func (r *InnerIsSysAdminH) Call(in *InnerIsSysAdminIn) (*InnerIsSysAdminOut, error) {
	once_InnerIsSysAdmin.Do(func() {
		local_InnerIsSysAdmin_EP = new(InnerIsSysAdminH).ProxySD()
	})
	//
	ep := local_InnerIsSysAdmin_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return &InnerIsSysAdminOut{
			OK:     false,
			ErrStr: err.Error(),
		}, nil
	}

	return result.(*InnerIsSysAdminOut), nil
}
