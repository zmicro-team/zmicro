package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/examples/hello/proto"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zmicro.New(zmicro.WithInitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, args *proto.HelloRequest, reply *proto.HelloReply) (err error) {
	*reply = proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", args.Name),
	}

	return nil
}
