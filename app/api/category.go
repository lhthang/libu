package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/constant"
	"net/http"
)

func ApplyCategoryAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	categoryEntity := repository.NewCategoryEntity(resource)
	categoryRoute := app.Group("/categories")

	categoryRoute.GET("", getAllCategories(categoryEntity))
	categoryRoute.GET("/:id", getCategoryById(categoryEntity))

	categoryRoute.Use(middlewares.RequireAuthenticated())
	categoryRoute.Use(middlewares.RequireAuthorization(constant.ADMIN))
	categoryRoute.POST("", createCategory(categoryEntity))
	categoryRoute.PUT("/:id", updateCategory(categoryEntity))
	categoryRoute.DELETE("/:id", deleteCategory(categoryEntity))
}

// GetAllCategories godoc
// @Tags CategoryController
// @Summary Get all categories
// @Description Get all categories
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Category
// @Router /categories [get]
func getAllCategories(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		categories, code, err := categoryEntity.GetAll()

		if err != nil {
			logrus.Print(err)
		}
		response := map[string]interface{}{
			"categories": categories,
			"err":        err.Error(),
		}
		ctx.JSON(code, response)
	}
}

// GetCategoryByID godoc
// @Tags CategoryController
// @Summary Get category by id
// @Description Get category by id
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {array} model.Category
// @Router /categories/{id} [get]
func getCategoryById(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		category, code, err := categoryEntity.GetOneByID(id)

		response := map[string]interface{}{
			"category": category,
			"err":      err.Error(),
		}
		ctx.JSON(code, response)
	}
}

// CreateCategory godoc
// @Tags CategoryController
// @Summary Create category
// @Description Create category
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param category body form.CategoryForm true "Category"
// @Success 200 {object} model.Category
// @Router /categories [post]
func createCategory(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		categoryForm := form.CategoryForm{}
		if err := ctx.Bind(&categoryForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		category, code, err := categoryEntity.CreateOne(categoryForm)

		response := map[string]interface{}{
			"category": category,
			"err":      err.Error(),
		}
		ctx.JSON(code, response)
	}
}

// UpdateCategory godoc
// @Tags CategoryController
// @Summary Update category
// @Description Update category
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Category ID"
// @Param category body form.CategoryForm true "CategoryForm"
// @Success 200 {object} model.Category
// @Router /categories/{id} [put]
func updateCategory(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		categoryForm := form.CategoryForm{}
		if err := ctx.Bind(&categoryForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		category, code, err := categoryEntity.Update(id, categoryForm)
		response := map[string]interface{}{
			"category": category,
			"err":      err.Error(),
		}
		ctx.JSON(code, response)
	}
}

// DeleteCategory godoc
// @Tags CategoryController
// @Summary Delete category
// @Description Delete category
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Category ID"
// @Success 200 {object} model.Category
// @Router /categories/{id} [delete]
func deleteCategory(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		category, code, err := categoryEntity.Delete(id)
		response := map[string]interface{}{
			"category": category,
			"err":      err.Error(),
		}
		ctx.JSON(code, response)
	}
}


