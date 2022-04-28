package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/core/log"
	zgin "github.com/iobrother/zmicro/core/transport/http"
	"github.com/iobrother/zmicro/examples/gin/swagger/api"
)

// curl http://127.0.0.1:5188/hello/zmicro
// http://127.0.0.1:5188/swagger/index.html
func main() {
	app := zmicro.New(zmicro.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitHttpServer(r *gin.Engine) error {
	Swagger(r)
	gin.DisableBindValidation()

	g := r.Group("/")
	proto.RegisterGreeterHTTPServer(g, &GreeterImpl{})

	return nil
}

type GreeterImpl struct {
	zgin.Implemented
}

func (s *GreeterImpl) SayHello(ctx context.Context, req *proto.HelloRequest) (rsp *proto.HelloReply, err error) {
	rsp = &proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return rsp, nil
}
