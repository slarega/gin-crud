package repository

import (
	"auth-service/helper"
	"auth-service/models"
	"errors"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Create(user models.User) models.User
	Update(form models.UserForm, userId int)
	Delete(userId int) (*models.User, error)
	FindAll() []models.User
	FindById(userId int) (*models.User, error)
	FindByEmail(email string) (models.User, bool)
}

type UserRepo struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) IUserRepo {
	return &UserRepo{
		Db: Db,
	}
}

func (repo UserRepo) Create(user models.User) models.User {
	result := repo.Db.Create(&user)
	helper.ErrorPanic(result.Error)
	return user
}

func (repo UserRepo) Update(form models.UserForm, userId int) {
	user, err := repo.FindById(userId)
	helper.ErrorPanic(err)

	user.Email = form.Email
	user.Password = form.Password

	result := repo.Db.Save(&user)
	helper.ErrorPanic(result.Error)
}

func (repo UserRepo) Delete(userId int) (*models.User, error) {
	user := new(models.User)
	if err := repo.Db.Where("guid = ?", userId).Delete(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo UserRepo) FindAll() []models.User {
	var users []models.User
	result := repo.Db.Find(&users)
	helper.ErrorPanic(result.Error)
	return users
}

func (repo UserRepo) FindById(userId int) (*models.User, error) {
	var user *models.User
	result := repo.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

func (repo UserRepo) FindByEmail(email string) (models.User, bool) {
	var user models.User

	if result := repo.Db.Where("email = ?", email).First(&user); result.Error != nil {
		return models.User{}, true
	} else {
		return user, false
	}
}
