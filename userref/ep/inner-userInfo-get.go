package ep

//last update date 2020-04-21
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

	user_ref "github.com/vhaoran/yiintf/userref"
)

const (
	InnerUserInfoGet_H_PATH = "/InnerUserInfoGet"
)

// postman
type (
	InnerUserInfoGetService interface {
		Exec(ctx context.Context, in *InnerUserInfoGetIn) (*ykit.Result, error)
	}

	//input  data
	InnerUserInfoGetIn struct {
		UID int64 `json:"uid"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Data interface{} `json:"data"`
	//}

	// handler implements
	InnerUserInfoGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerUserInfoGetH) MakeLocalEndpoint(svc InnerUserInfoGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  InnerUserInfoGet ###########")
		spew.Dump(ctx)

		in := request.(*InnerUserInfoGetIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *InnerUserInfoGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerUserInfoGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerUserInfoGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
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

	before := tran.ServerBefore(ykit.Jwt2ctx())
	opts := make([]tran.ServerOption, 0)
	opts = append(opts, before)
	opts = append(opts, options...)

	handler := tran.NewServer(
		ep,
		r.DecodeRequest,
		r.base.EncodeResponse,
		opts...)
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
var once_InnserUserInfoGet sync.Once
var local_InnserUserInfoGet_EP endpoint.Endpoint

func (r *InnerUserInfoGetH) Call(uid int64) (*user_ref.UserInfoRef, error) {
	once_InnserUserInfoGet.Do(func() {
		local_InnserUserInfoGet_EP = new(InnerUserInfoGetH).ProxySD()
	})
	//
	ep := local_InnserUserInfoGet_EP
	//
	in := &InnerUserInfoGetIn{
		UID: uid,
	}
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	i := result.(*ykit.Result)

	if i.Data != nil {
		return i.Data.(*user_ref.UserInfoRef), nil
	}

	return nil, err
}
