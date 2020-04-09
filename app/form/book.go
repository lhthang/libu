package form

type BookForm struct {
	ReleaseAt   string   `json:"releaseAt,omitempty"`
	Title       string   `json:"title" binding:"required"`
	Authors     []string `json:"authors,omitempty"`
	Publisher   string   `json:"publisher,omitempty"`
	Image       string   `json:"image"  binding:"required"`
	Description string   `json:"description,omitempty"`
	Link        string   `json:"link,omitempty"`
}
