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

	args := &proto.HelloRequest{
		Name: "zmicro",
	}

	reply, err := cli.SayHello(context.Background(), args)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("reply: %s", reply.Message)
}
