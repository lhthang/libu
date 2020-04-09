package api

import (
	"github.com/gin-gonic/gin"
	"libu/app/repository"
	"libu/my_db"
	err2 "libu/utils/err"
	"net/http"
)

func ApplyBookAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	bookEntity := repository.NewBookEntity(resource)

	bookRoute := app.Group("books")
	bookRoute.GET("", getAllBooks(bookEntity))

}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get all books
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Book
// @Router /books [get]
func getAllBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		books, code, err := entity.GetAll()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		response := map[string]interface{}{
			"books": books,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
