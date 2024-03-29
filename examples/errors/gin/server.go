package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zmicro-team/zmicro"
	"github.com/zmicro-team/zmicro/core/log"
	zgin "github.com/zmicro-team/zmicro/core/transport/http"
	"github.com/zmicro-team/zmicro/examples/errors/api"
	"github.com/zmicro-team/zmicro/examples/errors/errno"
)

// curl -w " status=%{http_code}" http://localhost:5180/error/internal
// curl -w " status=%{http_code}" http://localhost:5180/error/bad
// curl -w " status=%{http_code}" http://localhost:5180/error/biz
// curl -w " status=%{http_code}" http://localhost:5180/error/zmicro
func main() {
	app := zmicro.New(zmicro.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(r *gin.Engine) error {
	gin.DisableBindValidation()

	r.Use(zgin.CarrierInterceptor(zgin.NewCarry()))

	g := r.Group("/")
	api.RegisterGreeterHTTPServer(g, &GreeterImpl{})

	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) TestError(ctx context.Context, req *api.ErrorRequest, rsp *api.ErrorReply) error {
	if req.Name == "internal" {
		return errno.ErrInternalServer("服务器错误详请")
	} else if req.Name == "bad" {
		return errno.ErrBadRequest("请求参数错误详请")
	} else if req.Name == "biz" {
		return errno.ErrBizError()
	}

	*rsp = api.ErrorReply{
		Message: fmt.Sprintf("[%s] 不是错误", req.Name),
	}
	return nil
}

func (s *GreeterImpl) SayHello(ctx context.Context, req *api.HelloRequest, rsp *api.HelloReply) error {
	*rsp = api.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}
