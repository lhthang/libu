package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/api"
	"libu/db"
	"libu/middlewares"
	"os"
)

type Routes struct {
}

func (app Routes) StartGin() {
	r := gin.Default()
	publicRoute := r.Group("/api/v1")
	resource, err := db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any",middlewares.NewSwagger())
	r.NoRoute(func(context *gin.Context) {
		context.File("./template/route_not_found.html")
	})

	api.ApplyToDoAPI(publicRoute, resource)
	api.ApplyUserAPI(publicRoute, resource)
	r.Run(":"+os.Getenv("PORT"))
}
