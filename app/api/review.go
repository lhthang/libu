package api

import (
	"github.com/gin-gonic/gin"
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/jwt"
	"net/http"
)

func ApplyReviewAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	reviewEntity := repository.NewReviewEntity(resource)

	reviewRoute := app.Group("reviews")
	reviewRoute.Use(middlewares.RequireAuthenticated())
	reviewRoute.GET("", getAllReviews(reviewEntity))
	reviewRoute.POST("", createReview(reviewEntity))
	reviewRoute.PUT("/:id", updateReview(reviewEntity))
	reviewRoute.DELETE("/:id", deleteReview(reviewEntity))
}

// GetReviewById godoc
// @Tags ReviewController
// @Summary Get review by id
// @Description Get review by id
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Success 200 {object} model.Review
// @Router /reviews/{id} [get]
func getAllReviews(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		//username := jwt.GetUsername(ctx)
		review, code, err := reviewEntity.GetOneById(id)
		response := map[string]interface{}{
			"review": review,
			"error":  err,
		}
		ctx.JSON(code, response)
	}
}

// CreateReview godoc
// @Tags ReviewController
// @Summary Create review
// @Description Create review
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param review body form.ReviewForm true "Review"
// @Success 200 {object} model.Review
// @Router /reviews [post]
func createReview(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		reviewForm := form.ReviewForm{}

		if err := ctx.Bind(&reviewForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		username := jwt.GetUsername(ctx)

		reviewForm.Username = username
		review, code, err := reviewEntity.Create(reviewForm)
		response := map[string]interface{}{
			"review": review,
			"error":  err,
		}
		ctx.JSON(code, response)
	}
}

// UpdateReview godoc
// @Tags ReviewController
// @Summary Update review
// @Description Update review
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "ReviewID"
// @Param review body form.ReviewForm true "Review"
// @Success 200 {object} model.Review
// @Router /reviews/{id} [put]
func updateReview(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		reviewForm := form.ReviewForm{}

		if err := ctx.Bind(&reviewForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		username := jwt.GetUsername(ctx)

		reviewForm.Username = username
		review, code, err := reviewEntity.Update(id, username, reviewForm)
		response := map[string]interface{}{
			"review": review,
			"error":  err,
		}
		ctx.JSON(code, response)
	}
}

// DeleteReviewByID godoc
// @Tags ReviewController
// @Summary Delete review by id
// @Description  Delete review by id
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Success 200 {array} model.Review
// @Router /reviews/{id} [delete]
func deleteReview(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		username := jwt.GetUsername(ctx)

		review, code, err := reviewEntity.Delete(id, username)
		response := map[string]interface{}{
			"review": review,
			"error":  err,
		}
		ctx.JSON(code, response)
	}
}
