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
// auth: whr  date:2020/11/1216:48--------------------------
// ####请勿擅改此功能代码####
// 用途：获取bbs-prize分档价格上的价格
// 适用场景：
// 执行角色： levelPrice + offset
//---------------------------------------------

const (
	InnerBBBSPrizeLevelGet_H_PATH = "/InnerBBBSPrizeLevelGet"
)

//获取用户所有好友
type (
	InnerBBBSPrizeLevelGetService interface {
		Exec(in *InnerBBBSPrizeLevelGetIn) (*InnerBBBSPrizeLevelGetOut, error)
	}

	//input data
	InnerBBBSPrizeLevelGetIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	InnerBBBSPrizeLevelGetOut struct {
		Old    float64 `json:"old"`
		Offset float64 `json:"offset"`
		Price  float64 `json:"price"`
		//错误放这里
		ErrStr string `json:"err_str"`
	}

	// handler implements
	InnerBBBSPrizeLevelGetH struct {
		base ykit.RootTran
	}
)

func (r *InnerBBBSPrizeLevelGetH) MakeLocalEndpoint(svc InnerBBBSPrizeLevelGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBBBSPrizeLevelGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBBBSPrizeLevelGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBBBSPrizeLevelGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBBBSPrizeLevelGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBBBSPrizeLevelGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBBBSPrizeLevelGetH) HandlerLocal(service InnerBBBSPrizeLevelGetService,
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
func (r *InnerBBBSPrizeLevelGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBBBSPrizeLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBBBSPrizeLevelGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBBBSPrizeLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBBBSPrizeLevelGet sync.Once
var local_InnerBBBSPrizeLevelGet_EP endpoint.Endpoint

func (r *InnerBBBSPrizeLevelGetH) Call(in *InnerBBBSPrizeLevelGetIn) (*InnerBBBSPrizeLevelGetOut, error) {
	once_InnerBBBSPrizeLevelGet.Do(func() {
		local_InnerBBBSPrizeLevelGet_EP = new(InnerBBBSPrizeLevelGetH).ProxySD()
	})
	//
	ep := local_InnerBBBSPrizeLevelGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*InnerBBBSPrizeLevelGetOut), nil
}
