package form

import "libu/app/model"

type BookForm struct {
	ReleaseAt   string   `json:"releaseAt,omitempty"`
	Title       string   `json:"title" binding:"required"`
	CategoryIds []string `json:"categoryIds,omitempty"`
	Authors     []string `json:"authors,omitempty"`
	Publisher   string   `json:"publisher,omitempty"`
	Image       string   `json:"image"  binding:"required"`
	Description string   `json:"description,omitempty"`
	Link        string   `json:"link,omitempty"`
}

type BookResponse struct {
	*model.Book
	Reviews    []model.Review   `json:"reviews,omitempty"`
	Categories []model.Category `json:"categories,omitempty"`
}
