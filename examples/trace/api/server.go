package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmicro-team/zmicro"
	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/core/transport/rpc/client"
	"github.com/zmicro-team/zmicro/examples/proto"
)

// curl http://127.0.0.1:5180/hello/zmicro
func main() {
	app := zmicro.New(zmicro.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(r *gin.Engine) error {
	r.GET("/hello/:name", func(c *gin.Context) {
		cc, err := client.NewClient(
			client.WithServiceName("Greeter"),
			client.WithServiceAddr("127.0.0.1:5188"),
			client.Tracing(true),
		)
		if err != nil {
			log.Error(err)
			return
		}
		cli := proto.NewGreeterClient(cc.GetXClient())

		args := &proto.HelloRequest{
			Name: c.Param("name"),
		}

		log.Infof(args.Name)

		reply, err := cli.SayHello(c.Request.Context(), args)
		if err != nil {
			log.Error(err)
			return
		}

		c.String(http.StatusOK, reply.Message)
	})

	return nil
}
