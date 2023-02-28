package main

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consul_api "github.com/hashicorp/consul/api"
	"os"
	"time"
	user_info "user_client/api"
)

var (
	Name = "bbs.user.client"
)

func callRPC(client user_info.UserInfoClient) {
	reply, err := client.GetUserInfo(context.Background(), &user_info.GetUserInfoRequest{UserID: 1})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[grpc] GetUserInfo %v\n", reply)
}

func main() {
	consulCli, err := consul_api.NewClient(consul_api.DefaultConfig())

	if err != nil {
		panic(err)
	}

	r := consul.New(consulCli)
	selector.SetGlobalSelector(wrr.NewBuilder())
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///bbs.user.info"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	gClient := user_info.NewUserInfoClient(conn)
	ch := make(chan bool)

	go func(client user_info.UserInfoClient, done chan bool) {
		for {
			select {
			case <-done:
				return
			default:
				time.Sleep(time.Second)
				callRPC(client)
			}
		}
	}(gClient, ch)

	id, _ := os.Hostname()

	client, err := consul_api.NewClient(consul_api.DefaultConfig())

	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Registrar(consul.New(client)),
	)
	if err := app.Run(); err != nil {
		ch <- true
		log.Fatal(err)
	}
}
