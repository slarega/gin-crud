package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthData struct {
	TokenId      int `gorm:"type:int;primary_key"`
	RefreshToken string
	User         User `gorm:"foreignKey:GUID; references:TokenId"`
}

type ClientTokens struct {
	AccessToken  string `json:"access_token" example:"access"`
	RefreshToken string `json:"refresh_token" example:"refresh"`
}

type TokenPayload struct {
	ClientIp string
	UserId   int
}

type TokenClaims struct {
	TokenPayload
	jwt.RegisteredClaims
}
