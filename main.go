package main

import (
	"libu/app"
)


// @title UIT-Libray API
// @version 1.0
// @description This is a backend for uit-library
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email lhthang.1998@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8585
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main(){
	var server app.Routes
	server.StartGin()
}