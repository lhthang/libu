package form

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
	FullName    string `json:"fullName" binding:"required"`
	OldPassword string `json:"oldPassword,omitempty"`
	Password    string `json:"password,omitempty"`
}

type UpdateUser struct {
	Usernames []string `json:"usernames,omitempty"`
	Roles     []string `json:"roles"`
}

type UserResponse struct {
	Id          string   `json:"id"`
	Username    string   `bson:"username" json:"username"`
	FullName    string   `bson:"fullName" json:"fullName"`
	FavoriteIds []string `bson:"favoriteIds" json:"favoriteIds"`
	Roles       []string `json:"roles"`
}
