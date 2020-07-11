package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                 primitive.ObjectID `bson:"_id" json:"id"`
	Username           string             `bson:"username" json:"username"`
	FullName           string             `bson:"fullName" json:"fullName"`
	Password           string             `bson:"password" json:"password"`
	Roles              []string           `bson:"roles" json:"roles"`
	FavoriteIds        []string           `bson:"favoriteIds" json:"favoriteIds"`
	FavoriteCategoryId []string           `bson:"favoriteCategoryIds" json:"favoriteCategoryIds"`
}
