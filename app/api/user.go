package api

import (
	"github.com/gin-gonic/gin"
	"libu/app/form"
	"libu/app/repository"
	"libu/middlewares"
	"libu/my_db"
	"libu/utils/bcrypt"
	"libu/utils/constant"
	err2 "libu/utils/err"
	"libu/utils/jwt"
	"net/http"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *my_db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("")
	authRoute.POST("/login", login(userEntity))
	authRoute.POST("/sign-up", signUp(userEntity))

	userRoute := app.Group("/users")
	userRoute.GET("/get-all", getAllUSer(userEntity))
	userRoute.Use(middlewares.RequireAuthenticated())
	userRoute.PUT("/update/:username", updateUser(userEntity))

	// when need authentication
	userRoute.Use(middlewares.RequireAuthorization(constant.ADMIN)) // when need authorization
	userRoute.GET("", getAllUSer(userEntity))
	userRoute.PUT("/update-roles", updateRole(userEntity))
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
// @Summary Get all user
// @Description Get all user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.User
// @Router /user [get]
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

// UpdateUser godoc
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
