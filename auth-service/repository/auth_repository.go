package repository

import (
	"auth-service/helper"
	"auth-service/models"
	"gorm.io/gorm"
)

type IAuthRepo interface {
	CreateOrUpdate(GUID int, token string)
	Get(GUID int) models.AuthData
}

type AuthRepo struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepo {
	return &AuthRepo{
		Db: Db,
	}
}

func (repo AuthRepo) Get(GUID int) models.AuthData {
	var data models.AuthData
	result := repo.Db.Find(&data, GUID)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo AuthRepo) CreateOrUpdate(GUID int, token string) {
	data := models.AuthData{
		TokenId:      GUID,
		RefreshToken: token,
	}
	result := repo.Db.Save(&data)
	helper.ErrorPanic(result.Error)
}
