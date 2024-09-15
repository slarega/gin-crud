package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/helper"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/router"
	"auth-service/service"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

// @title 						CRUD API
// @version						1.0
// @description 				API на Go с использованием Gin и Gorm
// @host 						localhost:8888
// @BasePath 					/api
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env")
	}

	db := config.DatabaseConnection()
	validate := validator.New()
	dbUsersErr := db.Table("users").AutoMigrate(&models.User{})
	helper.ErrorPanic(dbUsersErr)

	dbAuthDataErr := db.Table("auth_data").AutoMigrate(&models.AuthData{})
	helper.ErrorPanic(dbAuthDataErr)

	userRepository := repository.NewUserRepository(db)
	authRepository := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepository, validate)
	userController := controller.NewUserController(userService)

	authService := service.NewAuthService(userRepository, authRepository, validate)
	authController := controller.NewAuthController(authService)

	routes := router.NewRouter(authController, userController)

	server := &http.Server{
		Addr:    ":8888",
		Handler: routes,
	}

	server_err := server.ListenAndServe()
	helper.ErrorPanic(server_err)
}
