package grpc

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
)

type ClientConfig struct {
	Name        string          `json:"name" yaml:"name"`                   // 服务名称
	Network     string          `json:"network" yaml:"network"`             // tcp unix
	Trace       bool            `json:"trace" yaml:"trace"`                 // 是否开启trace
	DialTimeOut time.Duration   `json:"dial_time_out" yaml:"dial_time_out"` // 拨通超时
	Timeout     time.Duration   `json:"timeout" yaml:"timeout"`             // keepalive.ClientParameters
	Etcd        clientv3.Config `json:"etcd" yaml:"etcd"`                   // etcd config
	ConfJson    string          `json:"conf_json" yaml:"conf_json"`         // 负载均衡等配置 例如`{"loadBalancingPolicy":"round_robin"}`

}

// scheme://authority/endpoint
func NewClient(conf *ClientConfig, target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	if conf.Trace {
		grpc.EnableTracing = true
	}

	client, err := clientv3.New(conf.Etcd)
	if err != nil {
		return nil, err
	}
	build, err := resolver.NewBuilder(client)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.DialTimeOut)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithResolvers(build),
		grpc.WithDefaultServiceConfig(conf.ConfJson),
	}

	conn, err := grpc.DialContext(ctx, target, append(opts, opt...)...)
	if err != nil {
		log.Printf("grpc lottery client dial err:%s", err)
		return nil, err
	}
	log.Printf("%s client grpc running on [%s]\n", conf.Name, target)
	return conn, nil
}
