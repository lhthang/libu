package form

import "libu/app/model"

type User struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateInformation struct {
	FullName      string `json:"fullName,omitempty"`
	OldPassword   string `json:"oldPassword,omitempty"`
	Password      string `json:"password,omitempty"`
	ProfileAvatar string `json:"profileAvatar"`
	DataLink      string `json:"dataLink"`
}

type UpdateUser struct {
	Usernames []string `json:"usernames,omitempty"`
	Roles     []string `json:"roles"`
}

type UserResponse struct {
	Id                  string           `json:"id"`
	Username            string           `json:"username"`
	FullName            string           `json:"fullName"`
	FavoriteIds         []string         `json:"favoriteIds"`
	FavoriteCategoryIds []string         `json:"favoriteCategoryIds"`
	Books               []BookResponse   `json:"favoriteBooks"`
	Categories          []model.Category `json:"favoriteCategories"`
	Roles               []string         `json:"roles"`
	ProfileAvatar       string           `json:"profileAvatar"`
	DataLink            string           `json:"dataLink"`
}

type FavoriteForm struct {
	FavoriteId         string `json:"favoriteId"`
	FavoriteCategoryId string `json:"favoriteCategoryId"`
	Action             string `json:"action" binding:"oneof=ADD REMOVE"`
}
