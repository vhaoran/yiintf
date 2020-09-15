module github.com/vhaoran/yiintf

go 1.14

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-kit/kit v0.10.0
	github.com/vhaoran/vchat v1.9.9
	go.mongodb.org/mongo-driver v1.2.0
)

replace github.com/vhaoran/vchat => ../vchat

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

// 处理etcd编译出错
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
