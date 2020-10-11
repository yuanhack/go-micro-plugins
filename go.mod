module github.com/macheal/go-micro-plugins

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/kr/pretty v0.1.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.26.0

)

replace github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
