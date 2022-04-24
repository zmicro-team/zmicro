package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro"
	"github.com/iobrother/zmicro/core/log"
)

func main() {
	app := zmicro.New(zmicro.WithInitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitHttpServer(r *gin.Engine) error {
	r.GET("/hello/:name", func(c *gin.Context) {
		c.Query("name")
		c.Writer.WriteString(fmt.Sprintf("hello %s!", c.Param("name")))
	})

	return nil
}