package main

import (
	"context"

	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/core/transport/rpc/client"
	"github.com/zmicro-team/zmicro/examples/proto"
)

func main() {
	c, err := client.NewClient(
		client.WithServiceName("Greeter"),
		client.BasePath("/zmicro"),
		client.EtcdAddr([]string{"127.0.0.1:2379"}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	cli := proto.NewGreeterClient(c.GetXClient())

	req := &proto.HelloRequest{
		Name: "zmicro",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
