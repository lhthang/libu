package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	ReleaseAt   time.Time          `bson:"releaseAt" json:"releaseAt"`
	Title       string             `bson:"title" json:"title"`
	CategoryIds []string           `bson:"categoryIds" json:"categoryIds,omitempty"`
	AuthorIds   []string           `bson:"authorIds" json:"authorIds,omitempty"`
	Publisher   string             `bson:"publisher" json:"publisher,omitempty"`
	Image       string             `bson:"image" json:"image"`
	Description string             `bson:"description" json:"description,omitempty"`
	Link        string             `bson:"link" json:"link,omitempty"`
	// Authors     []string           `bson:"authors" json:"authors,omitempty"`
}
