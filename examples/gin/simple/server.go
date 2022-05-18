package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmicro-team/zmicro"
	"github.com/zmicro-team/zmicro/core/log"
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
		c.String(http.StatusOK, fmt.Sprintf("hello %s!", c.Param("name")))
	})

	return nil
}
