package form

import "libu/app/model"

type ReviewForm struct {
	Id       string `json:"id"`
	Comment  string `json:"comment"`
	BookID   string `json:"bookId" binding:"required"`
	Username string `son:"username"`
	Rating   int    `json:"rating" binding:"oneof=1 2 3 4 5"`
}

type ReviewResponse struct {
	Reviews   []model.Review `json:"reviews"`
	AvgRating float32        `json:"avgRating"`
}
