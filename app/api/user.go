package api

import (
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/bcrypt"
	"libu/utils/constant"
	err2 "libu/utils/err"
	"libu/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("")
	authRoute.GET("/check-token/:token", checkToken())
	authRoute.POST("/login", login(userEntity))
	authRoute.POST("/sign-up", signUp(userEntity))

	userRoute := app.Group("/users")
	userRoute.GET("/get-all", getAllUSer(userEntity))
	userRoute.GET("/get/:id", getUserById(userEntity))
	userRoute.Use(middlewares.RequireAuthenticated())
	userRoute.PUT("/update/:username", updateUser(userEntity))
	userRoute.PUT("/favorites", addFavorite(userEntity))
	// when need authentication
	userRoute.Use(middlewares.RequireAuthorization(constant.ADMIN)) // when need authorization
	userRoute.GET("", getAllUSer(userEntity))
	userRoute.PUT("/update-roles", updateRole(userEntity))
}

func checkToken() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Param("token")
		isValid := middlewares.ValidateToken(token)
		response := map[string]interface{}{
			"isValid": isValid,
		}
		ctx.JSON(200, response)
	}
}

func login(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		userRequest := form.LoginUser{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		user, code, _ := userEntity.GetOneByUsername(userRequest.Username)

		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Wrong username or password"})
			return
		}
		token := middlewares.GenerateJWTToken(*user)
		response := map[string]interface{}{
			"token": token,
			"error": nil,
		}
		ctx.JSON(code, response)
	}
}

// SignUp godoc
// @Tags UserController
// @Summary Sign up
// @Description Sign up
// @Accept  json
// @Produce  json
// @Param updateUser body form.User true "User"
// @Success 200 {object} model.User
// @Router /sign-up [post]
func signUp(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		userRequest := form.User{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		user, code, err := userEntity.CreateOne(userRequest)
		response := map[string]interface{}{
			"user":  user,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetAllUser godoc
// @Tags UserController
// @Summary Get all user
// @Description Get all user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} form.UserResponse
// @Router /users [get]
func getAllUSer(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, err := userEntity.GetAll()
		response := map[string]interface{}{
			"users": list,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetUserById godoc
// @Tags UserController
// @Summary Get user by Id
// @Description Get user by Id
// @Accept  json
// @Produce  json
// @Success 200 {object} form.UserResponse
// @Router /users/get/{id} [get]
func getUserById(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		user, code, err := userEntity.GetOneById(id)
		response := map[string]interface{}{
			"user":  user,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateUser godoc
// @Tags UserController
// @Summary Update user
// @Description Update user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Param updateUser body form.UpdateInformation true "Update User"
// @Success 200 {object} model.User
// @Router /users/update/{username} [put]
func updateUser(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		updateUserForm := form.UpdateInformation{}
		if err := ctx.Bind(&updateUserForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		username := ctx.Param("username")
		userRequest := jwt.GetUsername(ctx)
		if userRequest != username {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "can not update this user"})
			return
		}

		user, code, err := userEntity.UpdateUser(username, updateUserForm)
		response := map[string]interface{}{
			"users": user,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateRole godoc
// @Tags UserController
// @Summary Update users roles
// @Description Update users roles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param updateUser body form.UpdateUser true "Update User"
// @Success 200 {object} model.User
// @Router /users/update-roles [put]
func updateRole(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		updateUserForm := form.UpdateUser{}
		if err := ctx.Bind(&updateUserForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		user, code, errs := userEntity.UpdateRole(updateUserForm)
		response := map[string]interface{}{
			"users": user,
			"error": errs,
		}
		ctx.JSON(code, response)
	}
}

// AddFavorite godoc
// @Tags UserController
// @Summary Add favorite book
// @Description Add favorite book
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param favoriteForm body form.FavoriteForm true "Favorite Form"
// @Success 200 {object} form.UserResponse
// @Router /users/favorites [put]
func addFavorite(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		username := jwt.GetUsername(ctx)

		actionForm := form.FavoriteForm{}
		if err := ctx.Bind(&actionForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		user, code, err := userEntity.UpdateFavorite(actionForm, username)
		response := map[string]interface{}{
			"user":  user,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
