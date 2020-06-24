package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Report struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	CreateAt time.Time          `bson:"createAt" json:"createAt"`
	Reason   string             `bson:"reason" json:"reason"`
	Username string             `bson:"username" json:"username"`
	ReviewId string             `bson:"reviewId" json:"reviewId"`
}
