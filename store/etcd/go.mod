module github.com/macheal/go-plugins/store/etcd

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/kr/pretty v0.1.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.26.0

)

replace github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
