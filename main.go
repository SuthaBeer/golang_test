package main

import (
	"fmt"
	"golang-apiuser/middleware"

	AuthController "golang-apiuser/controller/auth"
	UserController "golang-apiuser/controller/user"
	DB "golang-apiuser/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "golang-apiuser/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 		GOLang API
// @version 	1.0
// @description API for testing
// @host 		localhost:3000

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	database_name := "golang_user"
	DB.InitDB(database_name)
	DB.InitDBColumns(database_name)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)

	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/allusers", UserController.GetAllUsers)
	authorized.GET("/userinfo", UserController.GetUserInfo)
	authorized.POST("/transfercredit", UserController.TransferCredit)
	authorized.GET("/transfercredithistory", UserController.TransferCreditHistory)
	r.Run("localhost:3000")
}
