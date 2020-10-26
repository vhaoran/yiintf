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
// auth: whr  date:2020/10/2611:24--------------------------
// ####请勿擅改此功能代码####
// 用途：：运营商闪断貼价格
// 适用场景：
// 执行角色：
//--------------------------------------------- 
const (
	InnerBLevelVie_H_PATH = "/InnerBLevelVie"
)

//获取用户所有好友
type (
	InnerBLevelVieService interface {
		Exec(in *InnerBLevelVieIn) (*InnerBLevelVieOut, error)
	}

	//input data
	InnerBLevelVieIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	InnerBLevelVieOut struct {
		BrokerID   int64 `json:"broker_id"`
		LevelID    int64 `json:"level_id"`
		LevelPrice int64 `json:"level_price"`
	}

	// handler implements
	InnerBLevelVieH struct {
		base ykit.RootTran
	}
)

func (r *InnerBLevelVieH) MakeLocalEndpoint(svc InnerBLevelVieService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBLevelVieIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBLevelVieH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBLevelVieIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBLevelVieH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBLevelVieOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBLevelVieH) HandlerLocal(service InnerBLevelVieService,
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
func (r *InnerBLevelVieH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBLevelVie_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBLevelVieH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBLevelVie_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBLevelVie sync.Once
var local_InnerBLevelVie_EP endpoint.Endpoint

func (r *InnerBLevelVieH) Call(in *InnerBLevelVieIn) (*InnerBLevelVieOut, error) {
	once_InnerBLevelVie.Do(func() {
		local_InnerBLevelVie_EP = new(InnerBLevelVieH).ProxySD()
	})
	//
	ep := local_InnerBLevelVie_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerBLevelVieOut), nil
}
