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
	BProductPriceGet_H_PATH = "/BProductPriceGet"
)

//获取用户所有好友
type (
	InnerBProductPriceGetService interface {
		Exec(in *InnerBProductPriceGetIn) (*InnerBProductPriceGetOut, error)
	}

	//input data
	InnerBProductPriceGetIn struct {
		BrokerID  int64  `json:"broker_id"`
		ProductID string `json:"product_id"`
		ColorCode string `json:"color_code"`
	}

	//output data
	InnerBProductPriceGetOut struct {
		Old    float64 `json:"old"`
		Offset float64 `json:"offset"`

		Price float64 `json:"price"`
	}

	// handler implements
	InnerBProductPriceGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerBProductPriceGetH) MakeLocalEndpoint(svc InnerBProductPriceGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBProductPriceGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBProductPriceGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBProductPriceGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBProductPriceGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBProductPriceGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBProductPriceGetH) HandlerLocal(service InnerBProductPriceGetService,
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
func (r *InnerBProductPriceGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		BProductPriceGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBProductPriceGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		BProductPriceGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBProductPriceGet sync.Once
var local_InnerBProductPriceGet_EP endpoint.Endpoint

func (r *InnerBProductPriceGetH) Call(in *InnerBProductPriceGetIn) (*InnerBProductPriceGetOut, error) {
	once_InnerBProductPriceGet.Do(func() {
		local_InnerBProductPriceGet_EP = new(InnerBProductPriceGetH).ProxySD()
	})
	//
	ep := local_InnerBProductPriceGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*InnerBProductPriceGetOut), nil
}
