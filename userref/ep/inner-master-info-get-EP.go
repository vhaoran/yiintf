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
	InnerMasterInfoGet_H_PATH = "/InnerMasterInfoGet"
)

// postman
type (
	InnerMasterInfoGetService interface {
		Exec(ctx context.Context, in *InnerMasterInfoGetIn) (*ykit.Result, error)
	}

	//input  data
	InnerMasterInfoGetIn struct {
		MasterID int64 `json:"master_id"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Body interface{} `json:"data"`
	//}

	// handler implements
	InnerMasterInfoGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerMasterInfoGetH) MakeLocalEndpoint(svc InnerMasterInfoGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  InnerMasterInfoGet ###########")
		spew.Dump(ctx)

		in := request.(*InnerMasterInfoGetIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *InnerMasterInfoGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerMasterInfoGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerMasterInfoGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerMasterInfoGetH) HandlerLocal(service InnerMasterInfoGetService,
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
func (r *InnerMasterInfoGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerMasterInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerMasterInfoGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerMasterInfoGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerMasterInfoGet sync.Once
var local_InnerMasterInfoGet_EP endpoint.Endpoint

func (r *InnerMasterInfoGetH) Call(in InnerMasterInfoGetIn) (*user_ref.MasterInfoRef, error) {
	once_InnerMasterInfoGet.Do(func() {
		local_InnerMasterInfoGet_EP = new(InnerMasterInfoGetH).ProxySD()
	})
	//
	ep := local_InnerMasterInfoGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	ret := result.(*ykit.Result)

	return ret.Data.(*user_ref.MasterInfoRef), nil
}
