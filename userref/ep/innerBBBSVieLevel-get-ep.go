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
	InnerBBBSVieLevelGet_H_PATH = "/InnerBBBSVieLevelGet"
)

//获取用户所有好友
type (
	InnerBBBSVieLevelGetService interface {
		Exec(in *InnerBBBSVieLevelGetIn) (*InnerBBBSVieLevelGetOut, error)
	}

	//input data
	InnerBBBSVieLevelGetIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	InnerBBBSVieLevelGetOut struct {
		Old    float64 `json:"old"`
		Offset float64 `json:"offset"`

		Price  float64 `json:"price"`
		//错误放这里

		ErrStr string  `json:"err_str"`
	}

	// handler implements
	InnerBBBSVieLevelGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerBBBSVieLevelGetH) MakeLocalEndpoint(svc InnerBBBSVieLevelGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBBBSVieLevelGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBBBSVieLevelGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBBBSVieLevelGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBBBSVieLevelGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBBBSVieLevelGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBBBSVieLevelGetH) HandlerLocal(service InnerBBBSVieLevelGetService,
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
func (r *InnerBBBSVieLevelGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBBBSVieLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBBBSVieLevelGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBBBSVieLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBBBSVieLevelGet sync.Once
var local_InnerBBBSVieLevelGet_EP endpoint.Endpoint

func (r *InnerBBBSVieLevelGetH) Call(in *InnerBBBSVieLevelGetIn) (*InnerBBBSVieLevelGetOut, error) {
	once_InnerBBBSVieLevelGet.Do(func() {
		local_InnerBBBSVieLevelGet_EP = new(InnerBBBSVieLevelGetH).ProxySD()
	})
	//
	ep := local_InnerBBBSVieLevelGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*InnerBBBSVieLevelGetOut), nil
}
