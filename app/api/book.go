package api

import (
	"github.com/gin-gonic/gin"
	"libu/app/form"
	"libu/app/repository"
	"libu/my_db"
	err2 "libu/utils/err"
	"net/http"
)

func ApplyBookAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	bookEntity := repository.NewBookEntity(resource)

	bookRoute := app.Group("books")
	bookRoute.GET("", getAllBooks(bookEntity))
	bookRoute.GET("/:id", getBookById(bookEntity))
	bookRoute.POST("", createBook(bookEntity))
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get all books
// @Accept  json
// @Produce  json
// @Param q query string false "Query"
// @Success 200 {array} form.BookResponse
// @Router /books [get]
func getAllBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		keyword:=ctx.Query("q")

		if keyword!=""{
			books, code, err := entity.Search(keyword)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
				return
			}
			response := map[string]interface{}{
				"books": books,
				"error": err2.GetErrorMessage(err),
			}
			ctx.JSON(code, response)
			return
		}

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

// CreateBook godoc
// @Summary Create book
// @Description Create book
// @Accept  json
// @Produce  json
// @Param bookForm body form.BookForm true "BookForm"
// @Success 200 {object} form.BookResponse
// @Router /books [post]
func createBook(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var bookForm form.BookForm

		if err := ctx.Bind(&bookForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		book, code, err := entity.Create(bookForm)

		response := map[string]interface{}{
			"book":  book,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetBookByID godoc
// @Summary Get book by id
// @Description Get book by id
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} form.BookResponse
// @Router /books/{id} [get]
func getBookById(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		book, code, err := entity.GetOneByID(id)

		response := map[string]interface{}{
			"book":  book,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
