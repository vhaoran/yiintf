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
	InnerApplyVerify_H_PATH = "/InnerApplyVerify"
)

//获取用户所有好友
type (
	InnerApplyVerifyService interface {
		Exec(in *InnerApplyVerifyIn) (*InnerApplyVerifyOut, error)
	}

	//input data
	InnerApplyVerifyIn struct {
		////master/broker
		//ApplyType string `json:"apply_type"`
		//uid
		UID int64 `json:"uid"`
	}

	//output data
	InnerApplyVerifyOut struct {
		Ok     bool   `json:"bool"`
		ErrStr string `json:"err_str"`
	}

	// handler implements
	InnerApplyVerifyH struct {
		base ykit.RootTran
	}
)

func (r *InnerApplyVerifyH) MakeLocalEndpoint(svc InnerApplyVerifyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerApplyVerifyIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerApplyVerifyH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerApplyVerifyIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerApplyVerifyH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerApplyVerifyOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerApplyVerifyH) HandlerLocal(service InnerApplyVerifyService,
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
func (r *InnerApplyVerifyH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerApplyVerify_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerApplyVerifyH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerApplyVerify_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerApplyVerify sync.Once
var local_InnerApplyVerify_EP endpoint.Endpoint

func (r *InnerApplyVerifyH) Call(in *InnerApplyVerifyIn) (*InnerApplyVerifyOut, error) {
	once_InnerApplyVerify.Do(func() {
		local_InnerApplyVerify_EP = new(InnerApplyVerifyH).ProxySD()
	})
	//
	ep := local_InnerApplyVerify_EP
	//
	result, _ := ep(context.Background(), in)

	return result.(*InnerApplyVerifyOut), nil
}
