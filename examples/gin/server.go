package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/iobrother/zmicro"
)

func main() {
	app := zmicro.New(zmicro.WithInitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(r *gin.Engine) error {
	r.GET("/hello/:name", func(c *gin.Context) {
		c.Query("name")
		c.Writer.WriteString(fmt.Sprintf("hello %s!", c.Param("name")))
	})

	return nil
}
