package form

import "libu/app/model"

type ReviewForm struct {
	Id       string  `json:"id"`
	Comment  string  `json:"comment"`
	BookID   string  `json:"bookId" binding:"required"`
	Username string  `son:"username"`
	Rating   float32 `json:"rating" binding:"min=0,max=5"`
	// Rating   int    `json:"rating" binding:"oneof=1 2 3 4 5"`
}

type ActionForm struct {
	Action string `json:"action" binding:"oneof= upvote unvote"`
}

type ReviewResp struct {
	*model.Review
	User        UserComment    `json:"user"`
	Reports     []model.Report `json:"reports"`
	UpvoteCount int            `json:"upvoteCount"`
	ReportCount int            `json:"reportCount"`
}

type ReviewResponse struct {
	ReviewResp []ReviewResp `json:"reviews"`
	AvgRating  float32      `json:"avgRating"`
}
