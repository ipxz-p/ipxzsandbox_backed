package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ipxsandbox/internal/handler/auth_handler"
	"github.com/ipxsandbox/internal/handler/user_handler"
	"github.com/ipxsandbox/internal/middleware"
	"github.com/ipxsandbox/internal/repository/user"
	authUsecase "github.com/ipxsandbox/internal/usecase/auth_usercase"
	userUsecase "github.com/ipxsandbox/internal/usecase/user"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := user.New(db)
	authUC := authUsecase.NewAuthUsecase(userRepo)
	userUC := userUsecase.NewUserUsecase(userRepo)

	authHandler := auth_handler.NewAuthHandler(authUC)
	userHandler := user_handler.NewUserHandler(userUC)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh-token", authHandler.RefreshToken)

	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	auth.GET("/users", userHandler.GetUsers)
	auth.POST("/users", userHandler.CreateUser)
}