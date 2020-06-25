package repository

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"net/http"
	"time"
)

var ReportEntity IReport

type reportEntity struct {
	db   *my_db.Resource
	repo *mongo.Collection
}

type IReport interface {
	CreateOne(reportForm form.ReportForm, username string) (model.Report, int, error)
	GetByReviewId(id string) ([]model.Report, int, error)
	GetReviewsByReportCount(reportCount int64) ([]primitive.ObjectID, int, error)
	DeleteByReviewId(id string) ([]model.Report, int, error)
}

func NewReportEntity(resource *my_db.Resource) IReport {
	reportRepo := resource.DB.Collection("report")
	ReportEntity = &reportEntity{
		db:   resource,
		repo: reportRepo,
	}
	return ReportEntity
}

func (entity *reportEntity) CreateOne(reportForm form.ReportForm, username string) (model.Report, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	review, _, err := ReviewEntity.GetOneById(reportForm.ReviewId)
	if review.Review == nil || err != nil {
		return model.Report{}, http.StatusNotFound, err
	}

	report := model.Report{
		Id:       primitive.NewObjectID(),
		CreateAt: time.Now(),
		Reason:   reportForm.Reason,
		Username: username,
		ReviewId: reportForm.ReviewId,
	}

	_, err = entity.repo.InsertOne(ctx, report)
	if err != nil {
		return model.Report{}, http.StatusBadRequest, err
	}
	return report, http.StatusOK, nil
}

func (entity *reportEntity) GetByReviewId(id string) ([]model.Report, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var reports []model.Report

	cursor, err := entity.repo.Find(ctx, bson.M{"reviewId": id})
	if err != nil {
		return reports, http.StatusBadRequest, err
	}

	for cursor.Next(ctx) {
		var report model.Report
		err := cursor.Decode(&report)
		if err != nil {
			logrus.Print(err)
		}
		reports = append(reports, report)
	}
	return reports, http.StatusOK, nil
}

func (entity *reportEntity) GetReviewsByReportCount(reportCount int64) ([]primitive.ObjectID, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var reviewIds []primitive.ObjectID

	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$reviewId", "total": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"total": -1}},
		{"$match": bson.M{"total": bson.M{"$gte": reportCount}}}}
	result, err := entity.repo.Aggregate(ctx, pipeline)

	if err != nil {
		return []primitive.ObjectID{}, http.StatusBadRequest, err
	}
	for result.Next(ctx) {
		var reviewGroup model.ReviewGroup
		err := result.Decode(&reviewGroup)
		if err != nil {
			logrus.Println(err)
			continue
		}
		objId, _ := primitive.ObjectIDFromHex(reviewGroup.Id)
		reviewIds = append(reviewIds, objId)
	}

	return reviewIds, http.StatusOK, nil
}

func (entity *reportEntity) DeleteByReviewId(id string) ([]model.Report, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var reports []model.Report

	_, err := entity.repo.DeleteMany(ctx, bson.M{"reviewId": id})
	if err != nil {
		return reports, http.StatusBadRequest, err
	}
	return reports, http.StatusOK, nil
}
