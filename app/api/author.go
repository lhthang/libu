package api

import (
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/constant"
	err2 "libu/utils/err"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ApplyAuthorAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	authorEntity := repository.NewAuthorEntity(resource)
	authorRoute := app.Group("/authors")

	authorRoute.GET("", getAllAuthors(authorEntity))
	authorRoute.GET("/:id", getAuthorById(authorEntity))

	authorRoute.Use(middlewares.RequireAuthenticated())
	authorRoute.Use(middlewares.RequireAuthorization(constant.ADMIN))
	authorRoute.POST("", createAuthor(authorEntity))
	authorRoute.PUT("/:id", updateAuthor(authorEntity))
	authorRoute.DELETE("/:id", deleteAuthor(authorEntity))
}

// GetAllAuthor godoc
// @Tags AuthorController
// @Summary Get all authors
// @Description Get all authors
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Author
// @Router /authors [get]
func getAllAuthors(authorEntity repository.IAuthor) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authors, code, err := authorEntity.GetAll()

		if err != nil {
			logrus.Print(err)
		}
		response := map[string]interface{}{
			"authors": authors,
			"err":     err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetAuthorByID godoc
// @Tags AuthorController
// @Summary Get author by id
// @Description Get author by id
// @Accept  json
// @Produce  json
// @Param id path string true "Author ID"
// @Success 200 {array} model.Author
// @Router /authors/{id} [get]
func getAuthorById(authorEntity repository.IAuthor) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		author, code, err := authorEntity.GetOneByID(id)

		response := map[string]interface{}{
			"author": author,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// CreateAuthor godoc
// @Tags AuthorController
// @Summary Create author
// @Description Create author
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param author body form.AuthorForm true "Author"
// @Success 200 {object} model.Author
// @Router /authors [post]
func createAuthor(authorEntity repository.IAuthor) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		authorForm := form.AuthorForm{}
		if err := ctx.Bind(&authorForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		author, code, err := authorEntity.CreateOne(authorForm)

		response := map[string]interface{}{
			"author": author,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateAuthor godoc
// @Tags AuthorController
// @Summary Update author
// @Description Update author
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Param author body form.AuthorForm true "AuthorForm"
// @Success 200 {object} model.Author
// @Router /authors/{id} [put]
func updateAuthor(authorEntity repository.IAuthor) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		authorForm := form.AuthorForm{}
		if err := ctx.Bind(&authorForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		author, code, err := authorEntity.Update(id, authorForm)
		response := map[string]interface{}{
			"author": author,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// DeleteAuthor godoc
// @Tags AuthorController
// @Summary Delete author
// @Description Delete author
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Success 200 {object} model.Author
// @Router /authors/{id} [delete]
func deleteAuthor(authorEntity repository.IAuthor) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		author, code, err := authorEntity.Delete(id)
		response := map[string]interface{}{
			"author": author,
			"err":    err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
