package router

import (
	"auth-service/controller"
	_ "auth-service/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter(authController *controller.AuthController, userController *controller.UserController) *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Api docs ")
	})

	baseRouter := router.Group("/api")

	authRouter := baseRouter.Group("/auth")
	authRouter.POST("/registration", authController.Registration)
	authRouter.POST("/login", authController.Login)
	authRouter.POST("/logout", authController.Logout)
	authRouter.POST("/refresh-token", authController.RefreshAccessToken)
	authRouter.POST("/acc-token", authController.GetAccTokenData)

	userRouter := baseRouter.Group("/user")
	userRouter.GET("", userController.FindAll)
	userRouter.GET("/:userId", userController.FindById)
	userRouter.PATCH("/:userId", userController.Update)
	userRouter.DELETE("/:userId", userController.Delete)

	return router
}
