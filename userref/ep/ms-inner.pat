



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
	$type$_H_PATH = "/GetUserFriendsInner"
)

//获取用户所有好友
type (
    $type$Service interface {
		Exec(in *$type$In) ([]*$type$Out, error)
	}

	//input data
	$type$In struct {

	}

	//output data
	$type$Out struct {
	}

	// handler implements
	$type$H struct {
		base ykit.RootTran
	}
)

    func (r *$type$H) MakeLocalEndpoint(svc $type$Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*$type$In)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *$type$H) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new($type$In), ctx, req)
}

//个人实现,参数不能修改
func (r *$type$H) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response []*$type$Out
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *$type$H) HandlerLocal(service $type$Service,
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
func (r *$type$H) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		MSTAG,
		"POST",
		GetUserFriendsInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *$type$H) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		MSTAG,
		"POST",
		GetUserFriendsInner_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_$type$ sync.Once
var local_$type$_EP endpoint.Endpoint

func (r *$type$H) Call(in *$type$In) ([]*$type$Out, error) {
	once_$type$.Do(func() {
		local_$type$_EP = new($type$H).ProxySD()
	})
	//
	ep := local_$type$_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.([]*$type$Out), nil
}

