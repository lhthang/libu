package form

type ReportForm struct {
	Reason   string `json:"reason"`
	ReviewId string `json:"reviewId" binding:"required"`
}
