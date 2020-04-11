package repository

import (
	"errors"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"net/http"
	"time"
)

var ReviewEntity IReview

type reviewEntity struct {
	db   *my_db.Resource
	repo *mongo.Collection
}

type IReview interface {
	GetOneById(id string) (*model.Review, int, error)
	Create(reviewForm form.ReviewForm) (model.Review, int, error)
	Update(id, username string, reviewForm form.ReviewForm) (model.Review, int, error)
	Delete(id, username string) (model.Review, int, error)
}

func NewReviewEntity(resource *my_db.Resource) IReview {
	reviewRepo := resource.DB.Collection("review")
	ReviewEntity = &reviewEntity{
		db:   resource,
		repo: reviewRepo,
	}
	return ReviewEntity
}

func (entity *reviewEntity) GetOneById(id string) (*model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var review model.Review
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, getHTTPCode(err), err
	}

	err = entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &review, http.StatusOK, nil
}

func (entity *reviewEntity) Create(reviewForm form.ReviewForm) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	book, _, err := BookEntity.GetOneByID(reviewForm.BookID)

	if book.Book == nil || err != nil {
		return model.Review{}, getHTTPCode(err), err
	}
	review := model.Review{
		Id:       primitive.NewObjectID(),
		UpdateAt: time.Now(),
		Comment:  reviewForm.Comment,
		BookID:   reviewForm.BookID,
		Username: reviewForm.Username,
		Rating:   reviewForm.Rating,
	}

	_, err = entity.repo.InsertOne(ctx, review)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	return review, http.StatusOK, nil
}

func (entity *reviewEntity) Update(id, username string, reviewForm form.ReviewForm) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	review, _, err := entity.GetOneById(id)
	if err != nil || review == nil {
		return model.Review{}, getHTTPCode(err), err
	}
	if username != review.Username {
		return model.Review{}, http.StatusBadRequest, errors.New("this is not your review")
	}

	err = copier.Copy(review, reviewForm)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	review.UpdateAt = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": review}, opts).Decode(&review)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}
	return *review, http.StatusOK, nil
}

func (entity *reviewEntity) Delete(id, username string) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	review, _, err := entity.GetOneById(id)
	if err != nil || review == nil {
		return model.Review{}, http.StatusNotFound, err
	}

	if review.Username != username {
		return model.Review{}, http.StatusBadRequest, errors.New("this is not your review")
	}

	err = entity.repo.FindOneAndDelete(ctx, bson.M{"_id": objID, "username": username}).Decode(&review)
	if err != nil {
		return model.Review{}, http.StatusBadRequest, err
	}

	return *review, http.StatusOK, nil
}
