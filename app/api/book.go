package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/constant"
	err2 "libu/utils/err"
	"libu/utils/firebase"
	"net/http"
	"strconv"
	"strings"
)

func ApplyBookAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	bookEntity := repository.NewBookEntity(resource)

	bookRoute := app.Group("books")
	bookRoute.GET("", getAllBooks(bookEntity))
	bookRoute.GET("/get-high-rated", getHighRatedBooks(bookEntity))
	bookRoute.GET("/get-latest", getNewBooks(bookEntity))
	bookRoute.GET("/get-popular", getPopularBooks(bookEntity))
	bookRoute.GET("/book/:id/", getBookById(bookEntity))
	bookRoute.GET("/book/:id/similar", getSimilarBooks(bookEntity))
	bookRoute.GET("/recommend", getRecommendBooks(bookEntity))
	bookRoute.Use(middlewares.RequireAuthenticated())
	bookRoute.Use(middlewares.RequireAuthorization(constant.ADMIN))
	bookRoute.POST("", createBook(bookEntity))
	bookRoute.PUT("/:id", updateBook(bookEntity))
	bookRoute.DELETE("/:id", deleteBook(bookEntity))
	bookRoute.POST("/upload", uploadFile())
}

// GetAllBooks godoc
// @Tags BookController
// @Summary Get all books
// @Description Get all books
// @Accept  json
// @Produce  json
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Param q query string false "Query"
// @Success 200 {array} form.BookResponse
// @Router /books [get]
func getAllBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		keyword := ctx.Query("q")

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100000
		}

		if keyword != "" {
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

		books, code, err := entity.GetAll(skip, limit)
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

// GetNewBooks godoc
// @Tags BookController
// @Summary Get new books
// @Description Get new books
// @Accept  json
// @Produce  json
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Success 200 {array} form.BookResponse
// @Router /books/get-latest [get]
func getNewBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100000
		}

		books, code, err := entity.GetNewBooks(skip, limit)
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

// GetPopularBooks godoc
// @Tags BookController
// @Summary Get popular books
// @Description Get popular books
// @Accept  json
// @Produce  json
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Success 200 {array} form.BookResponse
// @Router /books/get-popular [get]
func getPopularBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100000
		}

		books, code, err := entity.GetPopularBooks(skip, limit)
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

// GetPopularBooks godoc
// @Tags BookController
// @Summary Get high rated books
// @Description Get high rated books
// @Accept  json
// @Produce  json
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Success 200 {array} form.BookResponse
// @Router /books/get-high-rated [get]
func getHighRatedBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100000
		}

		books, code, err := entity.GetHighRatedBooks(skip, limit)
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

// GetPopularBooks godoc
// @Tags BookController
// @Summary Get popular books
// @Description Get popular books
// @Accept  json
// @Produce  json
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Param categories query string true "Categories"
// @Success 200 {array} form.BookResponse
// @Router /books/recommend [get]
func getRecommendBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 100000
		}

		categoryIds := strings.Split(ctx.Query("categories"), "*")

		books, code, err := entity.GetRecommendBooks(skip, limit, categoryIds)
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
// @Tags BookController
// @Summary Create book
// @Description Create book
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
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
// @Tags BookController
// @Summary Get book by id
// @Description Get book by id
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} form.BookResponse
// @Router /books/book/{id} [get]
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

// GetSimilarBooks godoc
// @Tags BookController
// @Summary Get similar books
// @Description Get similar books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param skip query int false "Skip"
// @Param limit query int false "Limit"
// @Success 200 {array} form.BookResponse
// @Router /books/book/{id}/similar [get]
func getSimilarBooks(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		skip, err := strconv.ParseInt(ctx.Query("skip"), 10, 64)
		if err != nil {
			skip = 0
		}
		limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			limit = 10
		}

		book, code, err := entity.GetSimilarBooks(skip, limit, id)

		response := map[string]interface{}{
			"book":  book,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateBook godoc
// @Tags BookController
// @Summary Update book by id
// @Description Update book by id
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param bookForm body form.UpdateBookForm true "BookForm"
// @Success 200 {object} form.BookResponse
// @Router /books/{id} [put]
func updateBook(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var bookForm form.UpdateBookForm

		if err := ctx.Bind(&bookForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		book, code, err := entity.Update(id, bookForm)
		response := map[string]interface{}{
			"book":  book,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// DeleteBook godoc
// @Tags BookController
// @Summary Delete book by id
// @Description Delete book by id
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Success 200 {object} form.BookResponse
// @Router /books/{id} [delete]
func deleteBook(entity repository.IBook) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		book, code, err := entity.Delete(id)

		response := map[string]interface{}{
			"book":  book,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UploadFile godoc
// @Tags BookController
// @Summary Upload file
// @Description Upload file
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param file formData file true "file"
// @Success 200 {object} string
// @Router /books/upload [post]
func uploadFile() func(ctx *gin.Context) {
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
