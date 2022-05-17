package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/core/log"
	zgin "github.com/iobrother/zmicro/core/transport/http"
	"github.com/iobrother/zmicro/examples/multi/server/api"
	"github.com/iobrother/zmicro/examples/proto"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zmicro.New(
		zmicro.InitRpcServer(InitRpcServer),
		zmicro.InitHttpServer(InitHttpServer),
	)

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

func (s *GreeterImpl) SayHello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloReply) error {
	*rsp = proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

func InitHttpServer(r *gin.Engine) error {
	gin.DisableBindValidation()

	g := r.Group("/")
	api.RegisterGreeterHTTPServer(g, &HttpGreeter{})

	return nil
}

type HttpGreeter struct {
	zgin.Implemented
}

func (s *HttpGreeter) SayHello(ctx context.Context, req *api.HelloRequest, rsp *api.HelloReply) error {
	*rsp = api.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}
