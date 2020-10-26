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
// auth: whr  date:2020/10/2611:22--------------------------
// ####请勿擅改此功能代码####
// 用途：运营商悬赏貼价格
// 适用场景：
// 执行角色：
//---------------------------------------------
const (
	InnerBLevelPrize_H_PATH = "/InnerBLevelPrize"
)

//获取用户所有好友
type (
	InnerBLevelPrizeService interface {
		Exec(in *InnerBLevelPrizeIn) (*InnerBLevelPrizeOut, error)
	}

	//input data
	InnerBLevelPrizeIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	InnerBLevelPrizeOut struct {
		BrokerID   int64 `json:"broker_id"`
		LevelID    int64 `json:"level_id"`
		LevelPrice int64 `json:"level_price"`
	}

	// handler implements
	InnerBLevelPrizeH struct {
		base ykit.RootTran
	}
)

func (r *InnerBLevelPrizeH) MakeLocalEndpoint(svc InnerBLevelPrizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBLevelPrizeIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBLevelPrizeH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBLevelPrizeIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBLevelPrizeH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBLevelPrizeOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBLevelPrizeH) HandlerLocal(service InnerBLevelPrizeService,
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
func (r *InnerBLevelPrizeH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBLevelPrize_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBLevelPrizeH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBLevelPrize_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBLevelPrize sync.Once
var local_InnerBLevelPrize_EP endpoint.Endpoint

func (r *InnerBLevelPrizeH) Call(in *InnerBLevelPrizeIn) (*InnerBLevelPrizeOut, error) {
	once_InnerBLevelPrize.Do(func() {
		local_InnerBLevelPrize_EP = new(InnerBLevelPrizeH).ProxySD()
	})
	//
	ep := local_InnerBLevelPrize_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerBLevelPrizeOut), nil
}
