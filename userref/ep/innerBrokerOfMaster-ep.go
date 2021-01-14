package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vhaoran/vchat/lib/ylog"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	InnerBrokersOfMaster_H_PATH = "/InnerBrokersOfMaster"
)

//获取用户所有好友
type (
	InnerBrokersOfMasterService interface {
		Exec(in *InnerBrokersOfMasterIn) (*InnerBrokersOfMasterOut, error)
	}

	//input data
	InnerBrokersOfMasterIn struct {
		MasterID int64 `json:"master_id"`
	}

	//output data
	InnerBrokersOfMasterOut struct {
		Items  []int64 `json:"items"`
		ErrStr string  `json:"err_str"`
	}

	// handler implements
	InnerBrokersOfMasterH struct {
		base ykit.RootTran
	}
)

func (r *InnerBrokersOfMasterH) MakeLocalEndpoint(svc InnerBrokersOfMasterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerBrokersOfMasterIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerBrokersOfMasterH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerBrokersOfMasterIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerBrokersOfMasterH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerBrokersOfMasterOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerBrokersOfMasterH) HandlerLocal(service InnerBrokersOfMasterService,
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
func (r *InnerBrokersOfMasterH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerBrokersOfMaster_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerBrokersOfMasterH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerBrokersOfMaster_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnerBrokersOfMaster sync.Once
var local_InnerBrokersOfMaster_EP endpoint.Endpoint

func (r *InnerBrokersOfMasterH) Call(in *InnerBrokersOfMasterIn) (*InnerBrokersOfMasterOut, error) {
	once_InnerBrokersOfMaster.Do(func() {
		local_InnerBrokersOfMaster_EP = new(InnerBrokersOfMasterH).ProxySD()
	})
	//
	ep := local_InnerBrokersOfMaster_EP
	//
	result, err := ep(context.Background(), in)
	if err != nil {
		err = errors.New(fmt.Sprint("错误：innerBrokerOfMaster-ep.go--- ", err.Error()))
		ylog.Error(err.Error())
		return nil, nil
	}

	if result != nil {
		return result.(*InnerBrokersOfMasterOut), nil
	}
	return nil, nil
}
