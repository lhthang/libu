package repository

import (
	"errors"
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"libu/utils/arrays"
	"libu/utils/bcrypt"
	"libu/utils/constant"
	"net/http"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserEntity IUser

type userEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type IUser interface {
	GetAll() ([]form.UserResponse, int, error)
	GetOneById(id string) (*form.UserResponse, int, error)
	GetOneByUsername(username string) (*model.User, int, error)
	CreateOne(userForm form.User) (*model.User, int, error)
	UpdateUser(username string, userForm form.UpdateInformation) (*form.UserResponse, int, error)
	UpdateRole(updateUser form.UpdateUser) ([]model.User, int, []string)
	UpdateFavorite(actionForm form.FavoriteForm, username string) (*form.UserResponse, int, error)
}

//func NewToDoEntity
func NewUserEntity(resource *my_db.Resource) IUser {
	userRepo := resource.DB.Collection("user")
	UserEntity = &userEntity{resource: resource, repo: userRepo}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]form.UserResponse, int, error) {
	usersList := []form.UserResponse{}
	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		logrus.Print(err)
		return []form.UserResponse{}, 400, err
	}

	for cursor.Next(ctx) {
		var user model.User
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		usersList = append(usersList, form.UserResponse{
			Id:                  user.Id.Hex(),
			Username:            user.Username,
			FullName:            user.FullName,
			FavoriteIds:         user.FavoriteIds,
			FavoriteCategoryIds: user.FavoriteCategoryId,
			Roles:               user.Roles,
		})
	}
	return usersList, http.StatusOK, nil
}

func (entity *userEntity) GetOneByUsername(username string) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var user model.User
	err :=
		entity.repo.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *userEntity) GetOneById(id string) (*form.UserResponse, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var user model.User
	objID, _ := primitive.ObjectIDFromHex(id)
	err :=
		entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	userResp := getUserResponse(&user)
	userResp.Books = getFavoriteBooks(&user)
	userResp.Categories =getFavoriteCategories(&user)
	return &userResp, http.StatusOK, nil
}

func (entity *userEntity) CreateOne(userForm form.User) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user := model.User{
		Id:                 primitive.NewObjectID(),
		Username:           userForm.Username,
		FullName:           userForm.FullName,
		Password:           bcrypt.HashPassword(userForm.Password),
		Roles:              []string{constant.ADMIN, constant.USER},
		FavoriteIds:        []string{},
		FavoriteCategoryId: []string{},
	}
	found, _, _ := entity.GetOneByUsername(user.Username)
	if found != nil {
		return nil, http.StatusBadRequest, errors.New("Username is taken")
	}
	_, err := entity.repo.InsertOne(ctx, user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *userEntity) UpdateUser(username string, userForm form.UpdateInformation) (*form.UserResponse, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user, _, err := entity.GetOneByUsername(username)

	if err != nil || user == nil {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	password := user.Password
	newPassword := userForm.Password
	isUpdatePw := false
	if userForm.Password != "" && userForm.OldPassword != "" {
		isUpdatePw = true
	}
	if isUpdatePw {
		if err = bcrypt.ComparePasswordAndHashedPassword(userForm.OldPassword, user.Password); err != nil {
			return nil, http.StatusBadRequest, errors.New("old password is wrong")
		}

		userForm.Password = bcrypt.HashPassword(newPassword)
		err = copier.Copy(user, userForm)
		if userForm.Password == "" {
			user.Password = password
		}
	} else {
		err = copier.Copy(user, userForm)
		user.Password = password
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"username": username}, bson.M{"$set": user}, opts).Decode(&user)

	userResp := getUserResponse(user)
	return &userResp, http.StatusOK, nil
}

func getUserResponse(user *model.User) form.UserResponse {
	userResp := form.UserResponse{
		Id:                  user.Id.Hex(),
		Username:            user.Username,
		FullName:            user.FullName,
		FavoriteIds:         user.FavoriteIds,
		FavoriteCategoryIds: user.FavoriteCategoryId,
		Roles:               user.Roles,
	}
	return userResp
}

func getFavoriteCategories(user *model.User) []model.Category {
	var categories []model.Category
	var wg sync.WaitGroup
	categoryResp := make(chan *model.Category)
	for i, _ := range user.FavoriteCategoryId {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chann := new(model.Category)
			category, _, err := CategoryEntity.GetOneByID(user.FavoriteCategoryId[i])
			if err != nil {
				logrus.Println(err)
				return
			}
			chann = category
			categoryResp <- chann
		}(i)
	}

	go func() {
		wg.Wait()
		close(categoryResp)
	}()
	for category := range categoryResp {
		categories = append(categories, *category)
	}
	logrus.Println(categories)
	return categories
}

func getFavoriteBooks(user *model.User) []form.BookResponse {
	var booksResp []form.BookResponse
	var wg sync.WaitGroup
	bookResp := make(chan *form.BookResponse)
	for i, _ := range user.FavoriteIds {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chann := new(form.BookResponse)
			book, _, err := BookEntity.GetOneByID(user.FavoriteIds[i])
			if err != nil {
				logrus.Println(err)
				return
			}
			chann = &book
			bookResp <- chann
		}(i)
	}

	go func() {
		wg.Wait()
		close(bookResp)
	}()
	for book := range bookResp {
		booksResp = append(booksResp, *book)
	}
	return booksResp
}

func (entity *userEntity) UpdateRole(userForm form.UpdateUser) ([]model.User, int, []string) {
	ctx, cancel := initContext()
	defer cancel()
	var users []model.User
	var errs []string

	for _, name := range userForm.Roles {
		if name != constant.ADMIN && name != constant.USER {
			return users, http.StatusBadRequest, append(errs, "one of roles is invalid")
		}
	}

	for _, username := range userForm.Usernames {
		user, _, err := entity.GetOneByUsername(username)

		if err != nil || user == nil {
			errs = append(errs, username+" not found")
			continue
		}

		isReturnNewDoc := options.After
		opts := &options.FindOneAndUpdateOptions{
			ReturnDocument: &isReturnNewDoc,
		}
		err = entity.repo.FindOneAndUpdate(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"roles": userForm.Roles}}, opts).Decode(&user)
		users = append(users, *user)
	}

	return users, http.StatusOK, errs
}

func (entity *userEntity) UpdateFavorite(favoriteForm form.FavoriteForm, username string) (*form.UserResponse, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user, _, err := entity.GetOneByUsername(username)

	if err != nil || user == nil {
		return nil, getHTTPCode(err), err
	}

	update := map[string]interface{}{}

	fields :=map[string]interface{}{}

	//prepareUpdate
	if favoriteForm.Action == constant.ADD{
		if favoriteForm.FavoriteId!=""{
			if arrays.Contains(user.FavoriteIds, favoriteForm.FavoriteId) {
				return nil, http.StatusBadRequest, errors.New("this book is already added to your favorite list")
			}
			fields["favoriteIds"] =favoriteForm.FavoriteId
		}
		if favoriteForm.FavoriteCategoryId!=""{
			if arrays.Contains(user.FavoriteCategoryId, favoriteForm.FavoriteCategoryId) {
				return nil, http.StatusBadRequest, errors.New("this category is already added to your favorite list")
			}
			fields["favoriteCategoryIds"] =favoriteForm.FavoriteCategoryId

		}
		update["$push"] = fields

	}
	if favoriteForm.Action == constant.REMOVE{
		if favoriteForm.FavoriteId!=""{
			if !arrays.Contains(user.FavoriteIds, favoriteForm.FavoriteId) {
				return nil, http.StatusBadRequest, errors.New("this book is not added to your favorite list")
			}
			fields["favoriteIds"] =favoriteForm.FavoriteId
		}
		if favoriteForm.FavoriteCategoryId!=""{
			if !arrays.Contains(user.FavoriteCategoryId, favoriteForm.FavoriteCategoryId) {
				return nil, http.StatusBadRequest, errors.New("this category is not added to your favorite list")
			}
			fields["favoriteCategoryIds"] =favoriteForm.FavoriteCategoryId
		}
		update["$pull"] = fields
	}


	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	logrus.Println(update)
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"username": username}, update, opts).Decode(&user)

	userResp := getUserResponse(user)
	userResp.Books = getFavoriteBooks(user)
	userResp.Categories = getFavoriteCategories(user)
	return &userResp, http.StatusOK, nil
}
