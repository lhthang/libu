package form

type AuthorForm struct {
	// Id       string `json:"id"`
	Name     string `json:"name" binding:"required"`
	About    string `json:"about"`
	PhotoURL string `json:"photoUrl"`
}

type AuthorResponse struct {
	Name     string `json:"name"`
	About    string `json:"about"`
	PhotoURL string `json:"photoUrl"`
}
