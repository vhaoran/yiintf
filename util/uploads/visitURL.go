package uploads

//for snippet用于标准返回值的微服务接口

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vhaoran/yiintf/util"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	VisitURL_H_PATH = "/VisitURL"
)

type (
	VisitURLService interface {
		Exec(ctx context.Context, in *VisitURLIn) (*ykit.Result, error)
	}

	//input data
	VisitURLIn struct {
		Key     string `json:"key"`
		Expired int    `json:"expired"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Body interface{} `json:"data"`
	//}

	// handler implements
	VisitURLH struct {
		base ykit.RootTran
	}
)

func (r *VisitURLH) MakeLocalEndpoint(svc VisitURLService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  VisitURL ###########")
		spew.Dump(ctx)

		in := request.(*VisitURLIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *VisitURLH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(VisitURLIn), ctx, req)
}

//个人实现,参数不能修改
func (r *VisitURLH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *VisitURLH) HandlerLocal(service VisitURLService,
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
func (r *VisitURLH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		util.MSTAG,
		"POST",
		VisitURL_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *VisitURLH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		util.MSTAG,
		"POST",
		VisitURL_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}
