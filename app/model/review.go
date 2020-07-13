package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	UpdateAt time.Time          `bson:"updateAt" json:"updateAt"`
	Comment  string             `bson:"comment" json:"comment"`
	BookId   string             `bson:"bookId" json:"bookId"`
	Username string             `bson:"username" json:"username"`
	Rating   float32            `bson:"rating" json:"rating"`
	Upvotes  []string           `bson:"upvotes" json:"upvotes"`
	// Rating    int                `bson:"rating" json:"rating"`
}

type ReviewGroup struct {
	Id           string `bson:"_id"`
	TotalReports int    `bson:"total"`
}
