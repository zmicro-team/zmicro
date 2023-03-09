package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zmicro-team/zmicro/core/transport/http"
)

func main() {
	srv := http.NewServer()
	srv.Use(func(c *gin.Context) {
		log.Println("Use")
	})

	srv.UseEx(func(c *http.Context) {
		log.Println("UseEx")
	})

	srv.GET("/foo", func(c *gin.Context) {
		c.String(200, "foo")
	})

	srv.GetEx("/bar", func(c *http.Context) {
		c.String(200, "bar")
	})

	_ = srv.Run(":5180")
}
