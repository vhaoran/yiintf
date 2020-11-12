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
// auth: whr  date:2020/11/1216:41--------------------------
// ####请勿擅改此功能代码####
// 用途：获取代理 上大师的服务项目的价格
// 适用场景：计算方法：平台大理由价格+代理的offset
// 执行角色：
//---------------------------------------------
const (
	InnerBMasterCatePriceGet_H_PATH = "/InnerBMasterCatePriceGet"
)

//获取用户所有好友
type (
	InnerBMasterCatePriceGetService interface {
		Exec(in *InnerBMasterCatePriceGetIn) (*InnerBMasterCatePriceGetOut, error)
	}

	//input data
	InnerBMasterCatePriceGetIn struct {
		BrokerID int64 `json:"broker_id"`
		MasterID int64 `json:"master_id"`
		CateID   int64 `json:"cate_id"`
	}

	//output data
	InnerBMasterCatePriceGetOut struct {
		Old    float64 `json:"old"`
		Offset float64 `json:"offset"`

		//大师服务项目在代理上的价格
		Price float64 `json:"price"`
	}

	// handler implements
	InnerBMasterCatePriceGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerBMasterCatePriceGetH) MakeLocalEndpoint(svc InnerBMasterCatePriceGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBMasterCatePriceGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBMasterCatePriceGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBMasterCatePriceGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBMasterCatePriceGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBMasterCatePriceGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBMasterCatePriceGetH) HandlerLocal(service InnerBMasterCatePriceGetService,
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
func (r *InnerBMasterCatePriceGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBMasterCatePriceGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBMasterCatePriceGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBMasterCatePriceGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBMasterCatePriceGet sync.Once
var local_InnerBMasterCatePriceGet_EP endpoint.Endpoint

func (r *InnerBMasterCatePriceGetH) Call(in *InnerBMasterCatePriceGetIn) (*InnerBMasterCatePriceGetOut, error) {
	once_InnerBMasterCatePriceGet.Do(func() {
		local_InnerBMasterCatePriceGet_EP = new(InnerBMasterCatePriceGetH).ProxySD()
	})
	//
	ep := local_InnerBMasterCatePriceGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*InnerBMasterCatePriceGetOut), nil
}
