package api

import (
	"github.com/gin-gonic/gin"
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


// CreateReport godoc
// @Tags ReportController
// @Summary Create report
// @Description Create report
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param report body form.ReportForm true "Report"
// @Success 200 {object} model.Review
// @Router /reports [post]
func createReport(reportEntity repository.IReport) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		reportForm := form.ReportForm{}
		if err := ctx.Bind(&reportForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		username := jwt.GetUsername(ctx)

		report, code, err := reportEntity.CreateOne(reportForm,username)
		response := map[string]interface{}{
			"report": report,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}