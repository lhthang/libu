package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	About    string             `bson:"about" json:"about"`
	PhotoURL string             `bson:"photoUrl" json: "photoUrl"`
}
