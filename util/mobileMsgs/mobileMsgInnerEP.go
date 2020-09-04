package mobileMsgs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vhaoran/yiintf/util"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	MobileMsgInner_H_PATH = "/MobileMsgInner"
)

type (
	MobileMsgInnerService interface {
		Exec(ctx context.Context, in *MobileMsgInnerIn) (*ykit.Result, error)
	}

	//input data
	MobileMsgInnerIn struct {
		//
		Action int    `json:"action,omitempty"`
		Mobile string `json:"mobile,omitempty"`
		Text   string `json:"text,omitempty"`
	}

	//output data

	// handler implements
	MobileMsgInnerH struct {
		base ykit.RootTran
	}
)

func (r *MobileMsgInnerH) MakeLocalEndpoint(svc MobileMsgInnerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  MobileMsgInner ###########")
		spew.Dump(ctx)

		in := request.(*MobileMsgInnerIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *MobileMsgInnerH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(MobileMsgInnerIn), ctx, req)
}

//个人实现,参数不能修改
func (r *MobileMsgInnerH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *MobileMsgInnerH) HandlerLocal(service MobileMsgInnerService,
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
func (r *MobileMsgInnerH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		util.MSTAG,
		"POST",
		MobileMsgInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *MobileMsgInnerH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		util.MSTAG,
		"POST",
		MobileMsgInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_MobileMsgInner sync.Once
var local_MobileMsgInner_EP endpoint.Endpoint

func (r *MobileMsgInnerH) Call(in *MobileMsgInnerIn) (*ykit.Result, error) {
	once_MobileMsgInner.Do(func() {
		local_MobileMsgInner_EP = new(MobileMsgInnerH).ProxySD()
	})
	//
	ep := local_MobileMsgInner_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*ykit.Result), nil
}
