package controller

import (
	"auth-service/data/response"
	"auth-service/helper"
	"auth-service/models"
	"auth-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	accessTime  = 5
	refreshTime = 60 * 5
)

type AuthController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Registration     godoc
// @Summary			Регистрация
// @Description		Регистрация пользователя
// @Tags			auth
// @Produce			application/json
// @Param			user body models.UserForm true "Данные пользователя"
// @Success			200 {object} response.UserOkResponse
// @Router			/auth/registration [post]
func (controller *AuthController) Registration(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	createUser := models.UserForm{}
	err := ctx.ShouldBindJSON(&createUser)
	helper.ErrorPanic(err)

	newUser, errR := controller.authService.Registration(createUser)
	if errR != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Data: errR.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{Data: newUser})
}

// Login  		    godoc
// @Summary			Авторизация
// @Description		Возвращает пару access_token и refresh_token в cookie и Response body
// @Tags			auth
// @Produce			application/json
// @Param			user body models.UserForm true "Данные пользователя"
// @Success			200 {object} models.ClientTokens
// @Router			/auth/login [post]
func (controller *AuthController) Login(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	user := models.UserForm{}
	err := ctx.ShouldBindJSON(&user)
	helper.ErrorPanic(err)

	checkUser, checkErr := controller.authService.Login(user)
	if checkErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Data: checkErr.Error()})
		return
	}

	payload := models.TokenPayload{ClientIp: ctx.ClientIP(), UserId: checkUser.GUID}
	clientToken := controller.authService.CreateTokens(payload)
	controller.authService.StoringTokensInCookie(ctx, clientToken)
	ctx.JSON(http.StatusOK, clientToken)
}

// RefreshAccessToken  	godoc
// @Summary				Обновление токенов
// @Description			Обновление пары access_token и refresh_token в cookie и Response body
// @Tags				auth
// @Produce				application/json
// @Success				200 {object} response.Response
// @Router				/auth/refresh-token [post]
func (controller *AuthController) RefreshAccessToken(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	clientToken := controller.authService.RefreshAccessToken(ctx)
	controller.authService.StoringTokensInCookie(ctx, clientToken)
	ctx.JSON(http.StatusOK, clientToken)
}

// GetAccTokenData  	godoc
// @Summary				Проверка access токена
// @Description			Получение данных access токена
// @Tags				auth
// @Produce				application/json
// @Param				tokens body models.ClientTokens true "Токены пользователя"
// @Success				200 {object} response.Response
// @Router				/auth/acc-token [post]
func (controller *AuthController) GetAccTokenData(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	data := controller.authService.GetAccTokenData(ctx)
	ctx.JSON(http.StatusOK, data)
}

// Logout  		    godoc
// @Summary			Выход из системы
// @Description		Удаляет пару access_token и refresh_token в cookie
// @Tags			auth
// @Produce			application/json
// @Param			user body models.UserForm true "Данные пользователя"
// @Success			200 {object} models.ClientTokens
// @Router			/auth/logout [post]
func (controller *AuthController) Logout(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.SetCookie("access_token", "", -1, "/api", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/api", "localhost", false, true)
	ctx.JSON(http.StatusOK, "Cookie удалены")
}
