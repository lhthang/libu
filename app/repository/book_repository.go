package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	GetAll() ([]model.Book, int, error)
}

func NewBookEntity(resource *my_db.Resource) IBook {
	bookRepo := resource.DB.Collection("book")
	BookEntity = &bookEntity{
		resource: resource,
		repo:     bookRepo,
	}
	return BookEntity
}

func (entity bookEntity) GetAll() ([]model.Book, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var books []model.Book
	cursor, err := entity.repo.Find(ctx, bson.M{})
	if err != nil {
		return books, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		books = append(books, book)
	}
	return books, http.StatusOK, nil
}
