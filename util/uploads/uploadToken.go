package uploads

//for snippet用于标准返回值的微服务接口

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vhaoran/yiintf/util"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/endpoint"
	tran "github.com/go-kit/kit/transport/http"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	UploadToken_H_PATH = "/UploadToken"
)

type (
	UploadTokenService interface {
		Exec(ctx context.Context, in *UploadTokenIn) (*ykit.Result, error)
	}

	//input data
	UploadTokenIn struct {
		//0 : default,1 hours other: N hour
		Expired int `json:"expired,omitempty"`
	}

	//output data
	//Result struct {
	//	Code int         `json:"code"`
	//	Msg  string      `json:"msg"`
	//	Data interface{} `json:"data"`
	//}

	// handler implements
	UploadTokenH struct {
		base ykit.RootTran
	}
)

func (r *UploadTokenH) MakeLocalEndpoint(svc UploadTokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  UploadToken ###########")
		spew.Dump(ctx)

		in := request.(*UploadTokenIn)
		return svc.Exec(ctx, in)
	}
}

//个人实现,参数不能修改
func (r *UploadTokenH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(UploadTokenIn), ctx, req)
}

//个人实现,参数不能修改
func (r *UploadTokenH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response ykit.Result
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *UploadTokenH) HandlerLocal(service UploadTokenService,
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
func (r *UploadTokenH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		util.MSTAG,
		"POST",
		UploadToken_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *UploadTokenH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		util.MSTAG,
		"POST",
		UploadToken_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}
