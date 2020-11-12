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
	BBBSPrizeLevelGet_H_PATH = "/BBBSPrizeLevelGet"
)

//获取用户所有好友
type (
	BBBSPrizeLevelGetService interface {
		Exec(in *BBBSPrizeLevelGetIn) (*BBBSPrizeLevelGetOut, error)
	}

	//input data
	BBBSPrizeLevelGetIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	BBBSPrizeLevelGetOut struct {
		Price float64 `json:"price"`
	}

	// handler implements
	BBBSPrizeLevelGetH struct {
		base ykit.RootTran
	}
)

func (r *BBBSPrizeLevelGetH) MakeLocalEndpoint(svc BBBSPrizeLevelGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*BBBSPrizeLevelGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *BBBSPrizeLevelGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(BBBSPrizeLevelGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *BBBSPrizeLevelGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *BBBSPrizeLevelGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *BBBSPrizeLevelGetH) HandlerLocal(service BBBSPrizeLevelGetService,
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
func (r *BBBSPrizeLevelGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		BBBSPrizeLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *BBBSPrizeLevelGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		BBBSPrizeLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_BBBSPrizeLevelGet sync.Once
var local_BBBSPrizeLevelGet_EP endpoint.Endpoint

func (r *BBBSPrizeLevelGetH) Call(in *BBBSPrizeLevelGetIn) (*BBBSPrizeLevelGetOut, error) {
	once_BBBSPrizeLevelGet.Do(func() {
		local_BBBSPrizeLevelGet_EP = new(BBBSPrizeLevelGetH).ProxySD()
	})
	//
	ep := local_BBBSPrizeLevelGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*BBBSPrizeLevelGetOut), nil
}
