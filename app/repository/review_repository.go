package repository

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"libu/utils/arrays"
	"net/http"
	"time"
)

var ReviewEntity IReview

type reviewEntity struct {
	db   *my_db.Resource
	repo *mongo.Collection
}

type IReview interface {
	GetAll(report int64) ([]form.ReviewResp, int, error)
	GetOneById(id string) (form.ReviewResp, int, error)
	GetByBookId(bookId string) (*form.ReviewResponse, int, error)
	Create(reviewForm form.ReviewForm) (model.Review, int, error)
	Update(id, username string, reviewForm form.ReviewForm) (model.Review, int, error)
	Delete(id, username string) (model.Review, int, error)
	DeleteByAdmin(id string) (model.Review, int, error)
	Upvote(id, username, action string) (model.Review, int, error)
}

func NewReviewEntity(resource *my_db.Resource) IReview {
	reviewRepo := resource.DB.Collection("review")
	ReviewEntity = &reviewEntity{
		db:   resource,
		repo: reviewRepo,
	}
	return ReviewEntity
}

func getReportsOfReview(review model.Review) []model.Report {
	reports, _, err := ReportEntity.GetByReviewId(review.Id.Hex())
	if err != nil {
		return []model.Report{}
	}
	return reports
}

func getUserOfReview(review model.Review) *form.UserComment {
	user, _, err := UserEntity.GetOneByUsername(review.Username)
	if err != nil {
		return nil
	}
	resp:=form.UserComment{
		Id:            user.Id.Hex(),
		Username:      user.Username,
		FullName:      user.FullName,
		ProfileAvatar: user.ProfileAvatar,
	}
	return &resp
}

func (entity *reviewEntity) GetAll(report int64) ([]form.ReviewResp, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var reviews []form.ReviewResp

	arrayIds, _, err := ReportEntity.GetReviewsByReportCount(report)
	if err != nil {
		return []form.ReviewResp{}, http.StatusBadRequest, err
	}

	logrus.Println(arrayIds)

	query := map[string]interface{}{}
	if len(arrayIds) > 0 {
		query["_id"] = bson.M{"$in": arrayIds}
	}

	cursor, err := entity.repo.Find(ctx, query)
	if err != nil {
		return []form.ReviewResp{}, http.StatusBadRequest, err
	}

	for cursor.Next(ctx) {
		var review model.Review
		err := cursor.Decode(&review)
		if err != nil {
			logrus.Println(err)
			continue
		}
		reports := getReportsOfReview(review)

		reviewResp := form.ReviewResp{
			Review:      &review,
			Reports:     reports,
			UpvoteCount: len(review.Upvotes),
			ReportCount: len(reports),
		}
		reviews = append(reviews, reviewResp)
	}

	return reviews, http.StatusOK, nil
}

func (entity *reviewEntity) GetOneById(id string) (form.ReviewResp, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var review model.Review
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return form.ReviewResp{}, getHTTPCode(err), err
	}

	err = entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		return form.ReviewResp{}, http.StatusNotFound, err
	}

	if err != nil {
		return form.ReviewResp{}, http.StatusBadRequest, err
	}

	reports := getReportsOfReview(review)

	reviewResp := form.ReviewResp{
		Review:      &review,
		Reports:     reports,
		UpvoteCount: len(review.Upvotes),
		ReportCount: len(reports),
	}
	return reviewResp, http.StatusOK, nil
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
		BookId:   reviewForm.BookID,
		Username: reviewForm.Username,
		Rating:   reviewForm.Rating,
		Upvotes:  []string{},
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
	if err != nil || review.Review == nil {
		return model.Review{}, getHTTPCode(err), err
	}
	if username != review.Username {
		return model.Review{}, http.StatusBadRequest, errors.New("this is not your review")
	}

	err = copier.Copy(review.Review, reviewForm)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	review.UpdateAt = time.Now()

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": review.Review}, opts).Decode(&review)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}
	return *review.Review, http.StatusOK, nil
}

func (entity *reviewEntity) Delete(id, username string) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	review, _, err := entity.GetOneById(id)
	if err != nil || review.Review == nil {
		return model.Review{}, http.StatusNotFound, err
	}

	if review.Username != username {
		return model.Review{}, http.StatusBadRequest, errors.New("this is not your review")
	}

	err = entity.repo.FindOneAndDelete(ctx, bson.M{"_id": objID, "username": username}).Decode(&review)
	if err != nil {
		return model.Review{}, http.StatusBadRequest, err
	}

	_, _, _ = ReportEntity.DeleteByReviewId(id)

	return *review.Review, http.StatusOK, nil
}

func (entity *reviewEntity) GetByBookId(bookId string) (*form.ReviewResponse, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var reviewResp form.ReviewResponse

	cursor, err := entity.repo.Find(ctx, bson.M{"bookId": bookId})
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	for cursor.Next(ctx) {
		var review model.Review
		err = cursor.Decode(&review)
		if err != nil {
			logrus.Print(err)
			continue
		}
		reports := getReportsOfReview(review)
		user :=getUserOfReview(review)
		resp := form.ReviewResp{
			Review:      &review,
			Reports:     reports,
			UpvoteCount: len(review.Upvotes),
			ReportCount: len(reports),
			User: *user,
		}
		reviewResp.ReviewResp = append(reviewResp.ReviewResp, resp)
	}
	reviewResp.AvgRating = calculateRating(reviewResp.ReviewResp)

	return &reviewResp, http.StatusOK, nil
}

func calculateRating(reviews []form.ReviewResp) float32 {
	if len(reviews) == 0 {
		return 0
	}
	sum := 0
	for _, review := range reviews {
		sum = sum + review.Rating
	}
	rating := float32(sum) / float32(len(reviews))
	return rating
}

func (entity *reviewEntity) Upvote(id, username, action string) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	review, _, err := entity.GetOneById(id)
	if err != nil || review.Review == nil {
		return model.Review{}, getHTTPCode(err), err
	}

	query := map[string]interface{}{}
	//query["$set"] = review.Review

	if action == "upvote" {
		if arrays.Contains(review.Upvotes, username) {
			return model.Review{}, http.StatusBadRequest, errors.New("you already upvoted this review")
		}

		query["$push"] = bson.M{"upvotes": username}
	} else {
		if !arrays.Contains(review.Upvotes, username) {
			return model.Review{}, http.StatusBadRequest, errors.New("you don't vote this review")
		}

		query["$pull"] = bson.M{"upvotes": username}
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, query, opts).Decode(&review.Review)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}
	return *review.Review, http.StatusOK, nil
}

func (entity *reviewEntity) DeleteByAdmin(id string) (model.Review, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Review{}, getHTTPCode(err), err
	}

	review, _, err := entity.GetOneById(id)
	if err != nil || review.Review == nil {
		return model.Review{}, http.StatusNotFound, err
	}


	err = entity.repo.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		return model.Review{}, http.StatusBadRequest, err
	}

	_, _, _ = ReportEntity.DeleteByReviewId(id)

	return *review.Review, http.StatusOK, nil
}
