package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Create(reviewForm form.ReviewForm) (model.Review, int, error)
}

func NewReviewEntity(resource *my_db.Resource) IReview {
	reviewRepo := resource.DB.Collection("review")
	ReviewEntity = &reviewEntity{
		db:   resource,
		repo: reviewRepo,
	}
	return ReviewEntity
}

func (entity *reviewEntity) Create(reviewForm form.ReviewForm) (model.Review, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	review := model.Review{
		Id:       primitive.NewObjectID(),
		CreateAt: time.Now(),
		Comment:  reviewForm.Comment,
		BookID:   reviewForm.BookID,
		Username: reviewForm.Username,
		Rating:   reviewForm.Rating,
	}

	_, err := entity.repo.InsertOne(ctx, review)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	return review, http.StatusOK, nil
}
