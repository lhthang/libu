package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/form"
	"libu/app/repository"
	"libu/my_db"
	err2 "libu/utils/err"
	"libu/utils/firebase"
	"net/http"
)

func ApplyCategoryAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	categoryEntity := repository.NewCategoryEntity(resource)
	categoryRoute := app.Group("/categories")

	categoryRoute.GET("", getAllCategories(categoryEntity))
	categoryRoute.GET("/:id", getCategoryById(categoryEntity))
	categoryRoute.POST("/upload", uploadFile(categoryEntity))
	categoryRoute.POST("", createCategory(categoryEntity))
	categoryRoute.PUT("/:id", updateCategory(categoryEntity))
	categoryRoute.DELETE("/:id", deleteCategory(categoryEntity))
}

// GetAllCategories godoc
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
			"err":        err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetCategoryByID godoc
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
			"err":      err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

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
			"err":      err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

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
			"err":      err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func deleteCategory(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		category, code, err := categoryEntity.Delete(id)
		response := map[string]interface{}{
			"category": category,
			"err":      err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

func uploadFile(categoryEntity repository.ICategory) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		file, err := ctx.FormFile("file")
		if err != nil {
			logrus.Print(err)
		}
		resp, code, err := firebase.UploadFile(*file)
		response := map[string]interface{}{
			"resp": resp,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
