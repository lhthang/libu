package middlewares

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"libu/docs"
)

func NewSwagger() gin.HandlerFunc {
	docs.SwaggerInfo.Title = "UIT Library API"
	docs.SwaggerInfo.Description = "Swagger for UIT Library API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "http://libu20.herokuapp.com"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
