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
	"strconv"
)

func ApplyReviewAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	reviewEntity := repository.NewReviewEntity(resource)

	reviewRoute := app.Group("reviews")
	reviewRoute.GET("", getAllReviews(reviewEntity))
	reviewRoute.GET("/:id", getReviewById(reviewEntity))
	reviewRoute.GET("/:id/book", getAllReviewsByBook(reviewEntity))
	reviewRoute.Use(middlewares.RequireAuthenticated())
	reviewRoute.POST("", createReview(reviewEntity))
	reviewRoute.POST("/:id/action", upvoteReview(reviewEntity))
	reviewRoute.PUT("/:id", updateReview(reviewEntity))
	reviewRoute.DELETE("/:id", deleteReview(reviewEntity))
}

// GetAllReviews godoc
// @Tags ReviewController
// @Summary Get all reviews
// @Description Get all reviews
// @Accept  json
// @Produce  json
// @Param report query int false "Report"
// @Success 200 {array} model.Review
// @Router /reviews [get]
func getAllReviews(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		report, err := strconv.ParseInt(ctx.Query("report"), 10, 64)
		if err != nil {
			report = 0
		}
		review, code, err := reviewEntity.GetAll(report)
		response := map[string]interface{}{
			"review": review,
			"error":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetReviewById godoc
// @Tags ReviewController
// @Summary Get review by id
// @Description Get review by id
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} model.Review
// @Router /reviews/{id} [get]
func getReviewById(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		//username := jwt.GetUsername(ctx)
		review, code, err := reviewEntity.GetOneById(id)
		response := map[string]interface{}{
			"review": review,
			"error":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetAllReviewsByBook godoc
// @Tags ReviewController
// @Summary Get reviews by book
// @Description Get reviews by book
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {array} form.ReviewResponse
// @Router /reviews/{id}/book [get]
func getAllReviewsByBook(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		//username := jwt.GetUsername(ctx)
		review, code, err := reviewEntity.GetByBookId(id)
		response := map[string]interface{}{
			"review": review,
			"error":  err2.GetErrorMessage(err),
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
			"error":  err2.GetErrorMessage(err),
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
			"error":  err2.GetErrorMessage(err),
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
			"error":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// VoteReview godoc
// @Tags ReviewController
// @Summary Vote review
// @Description Vote review
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Review ID"
// @Param action body form.ActionForm true "Action"
// @Success 200 {object} model.Review
// @Router /reviews/{id}/action [post]
func upvoteReview(reviewEntity repository.IReview) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id:=ctx.Param("id")
		username := jwt.GetUsername(ctx)

		actionForm := form.ActionForm{}
		if err := ctx.Bind(&actionForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		review, code, err := reviewEntity.Upvote(id,username,actionForm.Action)
		response := map[string]interface{}{
			"review": review,
			"error":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
