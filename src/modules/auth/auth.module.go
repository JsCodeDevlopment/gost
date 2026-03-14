package auth

import (
	"gost/src/config"
	"gost/src/modules/users"

	"github.com/gin-gonic/gin"
)

func InitModule(router *gin.RouterGroup) {
	userRepo := users.NewUserRepository(config.DB)
	userService := users.NewUserService(userRepo)

	authService := NewAuthService(userService)
	authController := NewAuthController(authService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", authController.Logout)
	}
}
