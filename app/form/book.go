package form

import (
	"libu/app/model"
	"mime/multipart"
)

type BookForm struct {
	ReleaseAt   string                `form:"releaseAt,omitempty"`
	Title       string                `form:"title" binding:"required"`
	CategoryIds []string              `form:"categoryIds,omitempty"`
	Authors     []string              `form:"authors,omitempty"`
	Publisher   string                `form:"publisher,omitempty"`
	Image       string                `form:"image"  binding:"required"`
	Description string                `form:"description,omitempty"`
	Link        string                `form:"link,omitempty"`
	File        *multipart.FileHeader `form:"file,omitempty"`
}

type BookResponse struct {
	*model.Book
	Reviews    []model.Review   `json:"reviews,omitempty"`
	Categories []model.Category `json:"categories,omitempty"`
}
