package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/api"
	"libu/middlewares"
	"libu/my_db"
	"os"
)

type Routes struct {
}

func (app Routes) StartGin() {
	r := gin.Default()
	publicRoute := r.Group("/api/v1")
	resource, err := my_db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", middlewares.NewSwagger())

	r.Static("/template", "./template")

	r.NoRoute(func(context *gin.Context) {
		context.File("./template/route_not_found.html")
	})

	api.ApplyCategoryAPI(publicRoute, resource)
	api.ApplyBookAPI(publicRoute, resource)
	api.ApplyUserAPI(publicRoute, resource)
	api.ApplyReviewAPI(publicRoute, resource)
	api.ApplyAuthorAPI(publicRoute, resource)
	api.ApplyReportAPI(publicRoute,resource)
	r.Run(":" + os.Getenv("PORT"))
}
