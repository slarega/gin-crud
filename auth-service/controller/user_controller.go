package controller

import (
	"auth-service/data/response"
	"auth-service/helper"
	"auth-service/models"
	"auth-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(authService service.IUserService) *UserController {
	return &UserController{
		userService: authService,
	}
}

// FindAll    		godoc
// @Summary			Получить всех пользователей
// @Description		Возвращает список пользователей
// @Tags			user
// @Produce			application/json
// @Success			200 {object} response.UsersOkResponse
// @Router			/user [get]
func (controller *UserController) FindAll(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	usersResponse := controller.userService.FindAll()
	ctx.JSON(http.StatusOK, response.Response{Data: usersResponse})
}

// FindById  		    godoc
// @Summary				Получить пользователя
// @Description			Возвращает пользователя по id
// @Tags				user
// @Produce				application/json
// @Param				userId path string true "User id"
// @Success				200 {object} response.UserOkResponse
// @Router				/user/{userId} [get]
func (controller *UserController) FindById(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.ErrorPanic(err)
	userResponse, dbErr := controller.userService.FindById(id)
	helper.ErrorPanic(dbErr)
	ctx.JSON(http.StatusOK, response.Response{Data: userResponse})
}

// Update  		    	godoc
// @Summary				Обновить данные пользователя
// @Description			Обновить данные пользователя по id
// @Tags				user
// @Produce				application/json
// @Param				userId path string true "User id"
// @Param				user body models.UserForm true "Данные пользователя"
// @Success				200 {object} response.UserOkResponse
// @Router				/user/{userId} [patch]
func (controller *UserController) Update(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	updateUser := models.UserForm{}
	err := ctx.ShouldBindJSON(&updateUser)
	helper.ErrorPanic(err)

	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.ErrorPanic(err)
	userResponse := controller.userService.Update(updateUser, id)

	ctx.JSON(http.StatusOK, response.Response{Data: userResponse})
}

// Delete  		    	godoc
// @Summary				Удалить пользователя
// @Description			Удалить пользователя по id
// @Tags				user
// @Produce				application/json
// @Param				userId path string true "User id"
// @Success				200 {object} response.UserOkResponse
// @Router				/user/{userId} [delete]
func (controller *UserController) Delete(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	helper.ErrorPanic(err)
	deletedUser := controller.userService.Delete(id)

	ctx.JSON(http.StatusOK, response.Response{Data: deletedUser})
}
