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
	"github.com/iGoogle-ink/gopay"
	"github.com/vhaoran/vchat/lib/ykit"
)

const (
	InnerPayNotify_H_PATH = "/InnerPayNotify"
)

//获取用户所有好友
type (
	InnerPayNotifyService interface {
		Exec(in *InnerPayNotifyIn) (*InnerPayNotifyOut, error)
	}

	//input data
	InnerPayNotifyIn struct {
		WxData gopay.BodyMap `json:"wx_data"`

		//id
		OutTradeNo string `json:"out_trade_no"`
		//支付类型
		AccType int `json:"acc_type"`
		//存根
		StubNo string `json:"stub_no"`
	}

	//output data
	InnerPayNotifyOut struct {
		OK     bool   `json:"ok"`
		ErrStr string `json:"err_str"`
	}

	// handler implements
	InnerPayNotifyH struct {
		base ykit.RootTran
	}
)

func (r *InnerPayNotifyH) MakeLocalEndpoint(svc InnerPayNotifyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("#############  GetUserFriendsInner ###########")
		spew.Dump(ctx)

		in := request.(*InnerPayNotifyIn)
		return svc.Exec(in)
	}
}

//个人实现,参数不能修改
func (r *InnerPayNotifyH) DecodeRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return r.base.DecodeRequest(new(InnerPayNotifyIn), ctx, req)
}

//个人实现,参数不能修改
func (r *InnerPayNotifyH) DecodeResponse(_ context.Context, res *http.Response) (interface{}, error) {
	var response *InnerPayNotifyOut
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//handler for router，微服务本地接口，
func (r *InnerPayNotifyH) HandlerLocal(service InnerPayNotifyService,
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
func (r *InnerPayNotifyH) HandlerSD(mid []endpoint.Middleware,
	options ...tran.ServerOption) *tran.Server {
	return r.base.HandlerSD(
		context.Background(),
		msTag,
		"POST",
		InnerPayNotify_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse,
		mid,
		options...)
}

func (r *InnerPayNotifyH) ProxySD() endpoint.Endpoint {
	return r.base.ProxyEndpointSD(
		context.Background(),
		msTag,
		"POST",
		InnerPayNotify_H_PATH,
		r.DecodeRequest,
		r.DecodeResponse)
}

//只用于内部调用 ，不从风头调用
var once_InnertPayNotify sync.Once
var local_InnertPayNotify_EP endpoint.Endpoint

func (r *InnerPayNotifyH) Call(in *InnerPayNotifyIn) (*InnerPayNotifyOut, error) {
	once_InnertPayNotify.Do(func() {
		local_InnertPayNotify_EP = new(InnerPayNotifyH).ProxySD()
	})
	//
	ep := local_InnertPayNotify_EP
	//
	result, err := ep(context.Background(), in)

	if err != nil {
		return nil, err
	}

	return result.(*InnerPayNotifyOut), nil
}
