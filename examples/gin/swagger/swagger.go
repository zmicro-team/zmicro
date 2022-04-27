package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro/examples/gin/swagger/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"
)

func Swagger(r gin.IRouter) {
	r.GET("/swagger/*any", SwaggerHandler())
}

func SwaggerHandler() gin.HandlerFunc {
	swag.Register(swag.Name, new(docs.Docs))
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
