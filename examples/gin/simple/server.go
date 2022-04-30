package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/core/log"
)

// curl http://127.0.0.1:5180/hello/zmicro
func main() {
	app := zmicro.New(zmicro.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitHttpServer(r *gin.Engine) error {
	r.GET("/hello/:name", func(c *gin.Context) {
		c.Writer.WriteString(fmt.Sprintf("hello %s!", c.Param("name")))
	})

	return nil
}
