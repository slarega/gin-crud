package service

import (
	"auth-service/helper"
	"auth-service/models"
	"auth-service/repository"
	"github.com/go-playground/validator/v10"
)

type IUserService interface {
	FindById(userId int) (*models.User, error)
	FindAll() []models.User
	Update(form models.UserForm, userId int) *models.User
	Delete(userId int) *models.User
}

type UserService struct {
	UserRepo repository.IUserRepo
	Validate *validator.Validate
}

func NewUserService(userRepo repository.IUserRepo, validate *validator.Validate) IUserService {
	return &UserService{
		UserRepo: userRepo,
		Validate: validate,
	}
}

func (serv *UserService) FindAll() []models.User {
	result := serv.UserRepo.FindAll()
	var users []models.User
	for _, value := range result {
		user := models.User{
			GUID:     value.GUID,
			Email:    value.Email,
			Password: value.Password,
		}
		users = append(users, user)
	}
	return users
}

func (serv *UserService) FindById(userId int) (*models.User, error) {
	user, err := serv.UserRepo.FindById(userId)
	helper.ErrorPanic(err)
	return user, err
}

func (serv *UserService) Update(form models.UserForm, userId int) *models.User {
	serv.UserRepo.Update(form, userId)

	user, err := serv.UserRepo.FindById(userId)
	helper.ErrorPanic(err)
	return user
}

func (serv *UserService) Delete(userId int) *models.User {
	user, err := serv.UserRepo.Delete(userId)
	helper.ErrorPanic(err)
	return user
}
