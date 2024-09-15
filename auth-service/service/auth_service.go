package service

import (
	emailsender "auth-service/email-sender"
	"auth-service/helper"
	"auth-service/models"
	"auth-service/repository"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"strings"
	"time"
)

type IAuthService interface {
	Registration(form models.UserForm) (models.User, error)
	Login(form models.UserForm) (models.User, error)
	CreateTokens(payload models.TokenPayload) models.ClientTokens
	StoringTokensInCookie(ctx *gin.Context, clientToken models.ClientTokens)
	RefreshAccessToken(ctx *gin.Context) models.ClientTokens
	GetAccTokenData(ctx *gin.Context) string
}

type AuthService struct {
	UserRepo repository.IUserRepo
	AuthRepo repository.IAuthRepo
	Validate *validator.Validate
}

func NewAuthService(
	userRepo repository.IUserRepo,
	authRepo repository.IAuthRepo,
	validate *validator.Validate) IAuthService {
	return &AuthService{
		UserRepo: userRepo,
		AuthRepo: authRepo,
		Validate: validate,
	}
}

func (serv *AuthService) Registration(form models.UserForm) (models.User, error) {
	user, check := serv.UserRepo.FindByEmail(form.Email)
	if !check {
		return user, errors.New("пользователь с таким email уже существует")
	}
	userModel := form.ToModel()
	result := serv.UserRepo.Create(userModel)
	return result, nil
}

func (serv *AuthService) Login(form models.UserForm) (models.User, error) {
	user, check := serv.UserRepo.FindByEmail(form.Email)
	if check {
		return models.User{}, errors.New("пользователя с таким email не существует")
	}
	if !(form.Password == user.Password) {
		return models.User{}, errors.New("некорректный пароль")
	}
	return user, nil
}

func (serv *AuthService) createAccToken(payload models.TokenPayload) string {
	claims := models.TokenClaims{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	helper.ErrorPanic(err)

	return tokenString
}

func (serv *AuthService) encodeRefToken(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (serv *AuthService) decodeRefToken(token string) string {
	data, err := base64.StdEncoding.DecodeString(token)
	helper.ErrorPanic(err)
	return strings.Split(string(data), `"`)[1]
}

func (serv *AuthService) hashingRefToken(token string) string {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	helper.ErrorPanic(err)
	return string(hashedToken)
}

func (serv *AuthService) verifyRefreshToken(hashedRefToken, token string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedRefToken), []byte(token))
	return err == nil
}

func (serv *AuthService) CreateTokens(payload models.TokenPayload) models.ClientTokens {
	accToken := serv.createAccToken(payload)

	tokenData := fmt.Sprintf(`"%s, %d"`, payload.ClientIp, payload.UserId)
	refToken := serv.encodeRefToken(tokenData)
	hashedRefToken := serv.hashingRefToken(refToken)
	serv.AuthRepo.CreateOrUpdate(payload.UserId, hashedRefToken)

	return models.ClientTokens{AccessToken: accToken, RefreshToken: refToken}
}

func (serv *AuthService) StoringTokensInCookie(ctx *gin.Context, clientToken models.ClientTokens) {
	eatAccess, err := strconv.Atoi(os.Getenv("EAT_ACCESS"))
	helper.ErrorPanic(err)
	eatRefresh, err1 := strconv.Atoi(os.Getenv("EAT_REFRESH"))
	helper.ErrorPanic(err1)

	ctx.SetCookie("access_token", clientToken.AccessToken,
		eatAccess, "/api", "localhost", false, true)
	ctx.SetCookie("refresh_token", clientToken.RefreshToken,
		eatRefresh, "/api", "localhost", false, true)
}

func (serv *AuthService) RefreshAccessToken(ctx *gin.Context) models.ClientTokens {
	currentIp := ctx.ClientIP()
	refToken, refErr := ctx.Cookie("refresh_token")
	helper.ErrorPanic(refErr)
	refTokenData := strings.Split(serv.decodeRefToken(refToken), `,`)
	userId, _ := strconv.Atoi(strings.TrimSpace(refTokenData[1]))

	// токен
	hashedRefToken := serv.AuthRepo.Get(userId).RefreshToken
	if !serv.verifyRefreshToken(hashedRefToken, refToken) {
		return models.ClientTokens{AccessToken: "", RefreshToken: ""}
	}

	// IP
	if refTokenData[0] != currentIp {
		user, err := serv.UserRepo.FindById(userId)
		helper.ErrorPanic(err)
		emailsender.EmailSender(user.Email)
	}

	return serv.CreateTokens(
		models.TokenPayload{ClientIp: currentIp, UserId: userId},
	)
}

func (serv *AuthService) GetAccTokenData(ctx *gin.Context) string {
	accToken, errAcc := ctx.Cookie("access_token")
	helper.ErrorPanic(errAcc)

	var tokenClaims models.TokenClaims
	token, err := jwt.ParseWithClaims(accToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if !token.Valid {
		fmt.Println("invalid token")
	}

	return fmt.Sprintf(`UserId: %d, ClientIp: %s`, tokenClaims.UserId, tokenClaims.ClientIp)
}
