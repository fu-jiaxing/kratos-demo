package main

import (
	"context"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	userinfo "user_client/api"
)

var (
	Name = "bbs.user.client"
)

var (
	_metricClientQPS = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "bbs_user_client",
		Subsystem:   "requests",
		Name:        "throughput",
		Help:        "throughput",
		ConstLabels: nil,
	}, []string{"kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(_metricClientQPS)
}

func callRPC(client userinfo.UserInfoClient) {
	reply, err := client.GetUserInfo(context.Background(), &userinfo.GetUserInfoRequest{UserID: 1})
	if err != nil {
		log.Errorf("[grpc] GetUserInfo failed: ", err)
	} else {
		log.Infof("[grpc] GetUserInfo %v\n", reply)
	}
}

func main() {
	consulCli, err := consulapi.NewClient(consulapi.DefaultConfig())

	if err != nil {
		panic(err)
	}

	r := consul.New(consulCli)
	selector.SetGlobalSelector(wrr.NewBuilder())
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///bbs.user.info"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metrics.Client(metrics.WithRequests(prom.NewCounter(_metricClientQPS))),
			circuitbreaker.Client(),
			// metrics.Client(metrics.WithRequests(prom.NewCounter(_metricClientQPS))),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	gClient := userinfo.NewUserInfoClient(conn)
	ch := make(chan bool)

	go func(client userinfo.UserInfoClient, done chan bool) {
		for {
			select {
			case <-done:
				return
			default:
				// time.Sleep(time.Duration(rand.Intn(1000)))
				callRPC(client)
			}
		}
	}(gClient, ch)

	httpSrv := http.NewServer(
		http.Address(":8100"),
	)
	httpSrv.Handle("/metrics", promhttp.Handler())

	id, _ := os.Hostname()

	client, err := consulapi.NewClient(consulapi.DefaultConfig())

	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Registrar(consul.New(client)),
		kratos.Server(httpSrv),
	)
	if err := app.Run(); err != nil {
		ch <- true
		log.Fatal(err)
	}
}
