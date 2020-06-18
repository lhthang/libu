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
	"libu/utils/bcrypt"
	"libu/utils/constant"
	"net/http"
)

var UserEntity IUser

type userEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type IUser interface {
	GetAll() ([]form.UserResponse, int, error)
	GetOneByUsername(username string) (*model.User, int, error)
	CreateOne(userForm form.User) (*model.User, int, error)
	UpdateUser(username string, userForm form.UpdateInformation) (model.User, int, error)
	UpdateRole(updateUser form.UpdateUser) ([]model.User, int, []string)
	UpdateFavorite(id, username, action string) (model.User, int, error)
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
			Id:          user.Id.Hex(),
			Username:    user.Username,
			FullName:    user.FullName,
			FavoriteIds: user.FavoriteIds,
			Roles: user.Roles,
		})
	}
	return usersList, http.StatusOK, nil
}

func (entity *userEntity) GetOneByUsername(username string) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var user model.User
	err :=
		entity.repo.FindOne(ctx, bson.M{"username": username}, ).Decode(&user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *userEntity) CreateOne(userForm form.User) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user := model.User{
		Id:          primitive.NewObjectID(),
		Username:    userForm.Username,
		FullName:    userForm.FullName,
		Password:    bcrypt.HashPassword(userForm.Password),
		Roles:       []string{constant.ADMIN, constant.USER},
		FavoriteIds: []string{},
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

func (entity *userEntity) UpdateUser(username string, userForm form.UpdateInformation) (model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user, _, err := entity.GetOneByUsername(username)

	if err != nil || user == nil {
		return *user, http.StatusNotFound, errors.New("not found")
	}

	password := user.Password
	newPassword := userForm.Password
	isUpdatePw := false
	if userForm.Password != "" && userForm.OldPassword != "" {
		isUpdatePw = true
	}
	if isUpdatePw {
		if err = bcrypt.ComparePasswordAndHashedPassword(userForm.OldPassword, user.Password); err != nil {
			return *user, http.StatusBadRequest, errors.New("old password is wrong")
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

	return *user, http.StatusOK, nil
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

func (entity *userEntity) UpdateFavorite(id, username, action string) (model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user, _, err := entity.GetOneByUsername(username)

	if err != nil || user == nil {
		return model.User{}, getHTTPCode(err), err
	}

	update := map[string]interface{}{}
	if action == constant.ADD {
		if arrays.Contains(user.FavoriteIds, id) {
			return model.User{}, http.StatusBadRequest, errors.New("this is already added to your favorite list")
		}
		update = bson.M{"$push": bson.M{"favoriteIds": id}}
	} else {
		if !arrays.Contains(user.FavoriteIds, id) {
			return model.User{}, http.StatusBadRequest, errors.New("this is not added to your favorite list")
		}
		update = bson.M{"$pull": bson.M{"favoriteIds": id}}
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"username": username}, update, opts).Decode(&user)

	return *user, http.StatusOK, nil
}
