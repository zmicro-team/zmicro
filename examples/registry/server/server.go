package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/core/log"
	"github.com/iobrother/zmicro/examples/proto"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zmicro.New(zmicro.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloReply) error {
	*rsp = proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	fmt.Println("kkkkk")
	return nil
}
