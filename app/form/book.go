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
	Rating     float32          `json:"rating"`
}

type UpdateBookForm struct {
	ReleaseAt   string   `json:"releaseAt,omitempty"`
	Title       string   `json:"title" binding:"required"`
	CategoryIds []string `json:"categoryIds,omitempty"`
	Authors     []string `json:"authors,omitempty"`
	Publisher   string   `json:"publisher,omitempty"`
	Image       string   `json:"image,omitempty"`
	Description string   `json:"description,omitempty"`
	Link        string   `json:"link,omitempty"`
}
