package api

import (
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/constant"
	err2 "libu/utils/err"
	"libu/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyReportAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	reportEntity := repository.NewReportEntity(resource)
	reportRoute := app.Group("/reports")
	reviewRoute := app.Group("/reviews")

	reportRoute.Use(middlewares.RequireAuthenticated())
	reportRoute.POST("", createReport(reportEntity))
	reviewRoute.Use(middlewares.RequireAuthorization(constant.ADMIN))
	reviewRoute.GET("/:id/reports", getByReviewId(reportEntity))
}

// GetByReviewId godoc
// @Tags ReportController
// @Summary Get report by review id
// @Description Get report by review id
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} model.Report
// @Router /reviews/{id}/reports [get]
func getByReviewId(reportEntity repository.IReport) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		reports, code, err := reportEntity.GetByReviewId(id)
		response := map[string]interface{}{
			"reports": reports,
			"error":   err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
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

		report, code, err := reportEntity.CreateOne(reportForm, username)
		response := map[string]interface{}{
			"report": report,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
