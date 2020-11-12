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
	BBBSVieLevelGet_H_PATH = "/BBBSVieLevelGet"
)

//获取用户所有好友
type (
	BBBSVieLevelGetService interface {
		Exec(in *BBBSVieLevelGetIn) (*BBBSVieLevelGetOut, error)
	}

	//input data
	BBBSVieLevelGetIn struct {
		BrokerID int64 `json:"broker_id"`
		LevelID  int64 `json:"level_id"`
	}

	//output data
	BBBSVieLevelGetOut struct {
		Price float64 `json:"price"`
	}

	// handler implements
	BBBSVieLevelGetH struct {
		base ykit.RootTran
	}
)

func (r *BBBSVieLevelGetH) MakeLocalEndpoint(svc BBBSVieLevelGetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*BBBSVieLevelGetIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *BBBSVieLevelGetH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(BBBSVieLevelGetIn), ctx, req)
}

//个人实现,参数不能修改
func (r *BBBSVieLevelGetH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *BBBSVieLevelGetOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *BBBSVieLevelGetH) HandlerLocal(service BBBSVieLevelGetService,
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
func (r *BBBSVieLevelGetH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		BBBSVieLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *BBBSVieLevelGetH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		BBBSVieLevelGet_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_BBBSVieLevelGet sync.Once
var local_BBBSVieLevelGet_EP endpoint.Endpoint

func (r *BBBSVieLevelGetH) Call(in *BBBSVieLevelGetIn) (*BBBSVieLevelGetOut, error) {
	once_BBBSVieLevelGet.Do(func() {
		local_BBBSVieLevelGet_EP = new(BBBSVieLevelGetH).ProxySD()
	})
	//
	ep := local_BBBSVieLevelGet_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, nil
	}

	return result.(*BBBSVieLevelGetOut), nil
}
