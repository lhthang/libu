package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	err2 "libu/utils/err"
	"libu/utils/jwt"
	"net/http"
)

func ApplyReportAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	reportEntity := repository.NewReportEntity(resource)
	reportRoute := app.Group("/reports")

	reportRoute.Use(middlewares.RequireAuthenticated())
	reportRoute.POST("", createReport(reportEntity))
}

func createReport(reportEntity repository.IReport) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		reportForm := form.ReportForm{}
		if err := ctx.Bind(&reportForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		username := jwt.GetUsername(ctx)

		report, code, err := reportEntity.CreateOne(reportForm,username)
		logrus.Print("heeeeee")
		response := map[string]interface{}{
			"report": report,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}