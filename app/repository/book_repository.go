package repository

import (
	"github.com/araddon/dateparse"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"net/http"
)

var BookEntity IBook

type bookEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type IBook interface {
	GetAll() ([]form.BookResponse, int, error)
	Search(keyword string) ([]form.BookResponse, int, error)
	Create(bookForm form.BookForm) (form.BookResponse, int, error)
	GetOneByID(id string) (form.BookResponse, int, error)
	Delete(id string) (form.BookResponse, int, error)
}

func NewBookEntity(resource *my_db.Resource) IBook {
	bookRepo := resource.DB.Collection("book")
	BookEntity = &bookEntity{
		resource: resource,
		repo:     bookRepo,
	}
	return BookEntity
}

func getCategoryOfBook(book model.Book) []model.Category {
	var categories []model.Category
	for _, id := range book.CategoryIds {
		category, _, err := CategoryEntity.GetOneByID(id)
		if err != nil || category == nil {
			continue
		}
		categories = append(categories, *category)
	}
	return categories
}

func (entity bookEntity) GetAll() ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse
	cursor, err := entity.repo.Find(ctx, bson.M{})
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		booksResp = append(booksResp, form.BookResponse{
			Book:       &book,
			Reviews:    nil,
			Categories: getCategoryOfBook(book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Create(bookForm form.BookForm) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	releaseAt, err := dateparse.ParseAny(bookForm.ReleaseAt)
	if err != nil {
		logrus.Println(err)
	}
	var categoryIds []string
	for _, id := range bookForm.CategoryIds {
		category, _, _ := CategoryEntity.GetOneByID(id)
		if category != nil {
			categoryIds = append(categoryIds, id)
		}
	}
	book := model.Book{
		Id:          primitive.NewObjectID(),
		ReleaseAt:   releaseAt,
		Title:       bookForm.Title,
		Authors:     bookForm.Authors,
		Publisher:   bookForm.Publisher,
		CategoryIds: categoryIds,
		Image:       bookForm.Image,
		Description: bookForm.Description,
		Link:        bookForm.Link,
	}

	_, err = entity.repo.InsertOne(ctx, book)
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    nil,
		Categories: getCategoryOfBook(book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) GetOneByID(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var book model.Book
	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    nil,
		Categories: getCategoryOfBook(book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) Search(keyword string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse

	title := map[string]interface{}{}
	description := map[string]interface{}{}
	author := map[string]interface{}{}
	publisher := map[string]interface{}{}
	title["title"] = bson.M{"$regex": keyword, "$options": "i"}
	description["description"] = bson.M{"$regex": keyword, "$options": "i"}
	author["authors"] = bson.M{"$regex": keyword, "$options": "i"}
	publisher["publisher"] = bson.M{"$regex": keyword, "$options": "i"}
	query := map[string]interface{}{}
	query["$or"] = []interface{}{title, description, author, publisher}

	cursor, err := entity.repo.Find(ctx, query)
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		booksResp = append(booksResp, form.BookResponse{
			Book:       &book,
			Reviews:    nil,
			Categories: getCategoryOfBook(book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Delete(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	bookResp, _, err := entity.GetOneByID(id)

	if err != nil || bookResp.Book == nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	return bookResp, http.StatusOK, nil
}
