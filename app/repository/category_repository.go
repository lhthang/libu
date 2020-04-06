package repository

import (
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"net/http"
)

var CategoryEntity ICategory

type categoryEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type ICategory interface {
	GetAll() ([]model.Category, int, error)
	GetOneByID(id string) (*model.Category, int, error)
	CreateOne(categoryForm form.CategoryForm) (model.Category, int, error)
	Update(id string, categoryForm form.CategoryForm) (model.Category, int, error)
	Delete(id string) (model.Category, int, error)
}

func NewCategoryEntity(resource *my_db.Resource) ICategory {
	categoryRepo := resource.DB.Collection("category")
	CategoryEntity = &categoryEntity{
		resource: resource,
		repo:     categoryRepo,
	}
	return CategoryEntity
}

func (entity *categoryEntity) GetAll() ([]model.Category, int, error) {
	ctx, cancel := initContext()
	var categories []model.Category
	defer cancel()

	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		return []model.Category{}, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var category model.Category
		err = cursor.Decode(&category)
		if err != nil {
			logrus.Print(err)
		}
		categories = append(categories, category)
	}
	return categories, http.StatusOK, nil
}

func (entity *categoryEntity) GetOneByID(id string) (*model.Category, int, error) {
	var category model.Category
	ctx, cancel := initContext()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&category)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &category, http.StatusOK, nil
}

func (entity *categoryEntity) CreateOne(categoryForm form.CategoryForm) (model.Category, int, error) {

	ctx, cancel := initContext()
	defer cancel()

	category := model.Category{
		Id:   primitive.NewObjectID(),
		Name: categoryForm.Name,
	}

	_, err := entity.repo.InsertOne(ctx, category)
	if err != nil {
		return model.Category{}, http.StatusNotFound, err
	}

	return category, http.StatusOK, nil
}

func (entity *categoryEntity) Update(id string, todoForm form.CategoryForm) (model.Category, int, error) {
	var category *model.Category
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)

	category, _, err := entity.GetOneByID(id)
	if err != nil {
		return model.Category{}, http.StatusNotFound, nil
	}

	err = copier.Copy(category, todoForm) // this is why we need return a pointer: to copy value
	if err != nil {
		logrus.Error(err)
		return model.Category{}, getHTTPCode(err), err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": category}, opts).Decode(&category)

	return *category, http.StatusOK, nil
}

func (entity *categoryEntity) Delete(id string) (model.Category, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	category, _, err := entity.GetOneByID(id)
	if err != nil || category == nil {
		return model.Category{}, getHTTPCode(err), err
	}

	objID, _ := primitive.ObjectIDFromHex(id)

	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return model.Category{}, getHTTPCode(err), err
	}

	return *category, http.StatusOK, nil
}
