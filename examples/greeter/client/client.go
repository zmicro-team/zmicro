package main

import (
	"context"

	"github.com/iobrother/zmicro/core/log"
	"github.com/iobrother/zmicro/core/transport/rpc/client"
	"github.com/iobrother/zmicro/examples/proto"
)

func main() {
	c := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	cli := proto.NewGreeterClient(c.GetXClient())

	req := &proto.HelloRequest{
		Name: "zmicro",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
