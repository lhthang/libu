package repository

import (
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var AuthorEntity IAuthor

type authorEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type IAuthor interface {
	GetAll() ([]model.Author, int, error)
	GetOneByID(id string) (*model.Author, int, error)
	GetBooksByAuthorId(authorId string) ([]form.BookResponse, int, error)
	CreateOne(authorForm form.AuthorForm) (model.Author, int, error)
	Update(id string, authorForm form.AuthorForm) (model.Author, int, error)
	Delete(id string) (model.Author, int, error)
}

func NewAuthorEntity(resource *my_db.Resource) IAuthor {
	authorRepo := resource.DB.Collection("author")
	AuthorEntity = &authorEntity{
		resource: resource,
		repo:     authorRepo,
	}
	return AuthorEntity
}

func (entity *authorEntity) GetAll() ([]model.Author, int, error) {
	ctx, cancel := initContext()
	var authors []model.Author
	defer cancel()

	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		return []model.Author{}, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var author model.Author
		err = cursor.Decode(&author)
		if err != nil {
			logrus.Print(err)
		}
		authors = append(authors, author)
	}
	return authors, http.StatusOK, nil
}

func (entity *authorEntity) GetOneByID(id string) (*model.Author, int, error) {
	var author model.Author
	ctx, cancel := initContext()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&author)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &author, http.StatusOK, nil
}

func (entity *authorEntity) GetBooksByAuthorId(authorId string) ([]form.BookResponse, int, error) {
	allBooks, _, err := BookEntity.GetAll(0,10000000)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	var booksByAuthor []form.BookResponse
	for _, book := range allBooks {
		for _, author := range book.Authors {
			if author.Id.Hex() == authorId {
				booksByAuthor = append(booksByAuthor, book)
				break
			}
		}
	}
	return booksByAuthor, http.StatusOK, nil
}

func (entity *authorEntity) CreateOne(authorForm form.AuthorForm) (model.Author, int, error) {

	ctx, cancel := initContext()
	defer cancel()

	author := model.Author{
		Id:       primitive.NewObjectID(),
		Name:     authorForm.Name,
		About:    authorForm.About,
		PhotoURL: authorForm.PhotoURL,
	}

	_, err := entity.repo.InsertOne(ctx, author)
	if err != nil {
		return model.Author{}, http.StatusNotFound, err
	}

	return author, http.StatusOK, nil
}

func (entity *authorEntity) Update(id string, todoForm form.AuthorForm) (model.Author, int, error) {
	var author *model.Author
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)

	author, _, err := entity.GetOneByID(id)
	if err != nil {
		return model.Author{}, http.StatusNotFound, nil
	}

	err = copier.Copy(author, todoForm)
	if err != nil {
		logrus.Error(err)
		return model.Author{}, getHTTPCode(err), err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": author}, opts).Decode(&author)

	return *author, http.StatusOK, nil
}

func (entity *authorEntity) Delete(id string) (model.Author, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	author, _, err := entity.GetOneByID(id)
	if err != nil || author == nil {
		return model.Author{}, getHTTPCode(err), err
	}

	objID, _ := primitive.ObjectIDFromHex(id)

	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return model.Author{}, getHTTPCode(err), err
	}

	return *author, http.StatusOK, nil
}
