package form

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReviewForm struct {
	Id       primitive.ObjectID `json:"id"`
	Comment  string             `json:"comment"`
	BookID   string             `json:"bookId" binding:"required"`
	Username string             `son:"username"`
	Rating   int                `json:"rating" binding:"oneof=1 2 3 4 5"`
}
