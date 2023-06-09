package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smallnest/rpcx/server"
	"github.com/zmicro-team/zmicro"
	"github.com/zmicro-team/zmicro/core/log"
	zgin "github.com/zmicro-team/zmicro/core/transport/http"
	"github.com/zmicro-team/zmicro/examples/multi/server/api"
	"github.com/zmicro-team/zmicro/examples/proto"
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

	r.Use(zgin.CarrierInterceptor(zgin.NewCarry()))

	g := r.Group("/")
	api.RegisterGreeterHTTPServer(g, &HttpGreeter{})

	return nil
}

type HttpGreeter struct{}

func (s *HttpGreeter) SayHello(ctx context.Context, req *api.HelloRequest, rsp *api.HelloReply) error {
	*rsp = api.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}
