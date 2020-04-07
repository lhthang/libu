package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Review struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	UpdateAt time.Time          `bson:"updateAt" json:"updateAt"`
	Comment  string             `bson:"comment" json:"comment"`
	BookID   string             `bson:"bookId" json:"bookId"`
	Username string             `bson:"username" json:"username"`
	Rating   int                `bson:"rating" json:"rating"`
}
