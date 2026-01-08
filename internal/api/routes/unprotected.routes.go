package routes

import (
	"github.com/SecureParadise/go_attendence/internal/api/handlers"
	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/SecureParadise/go_attendence/internal/config"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupUnProtectedRoutes(router *gin.Engine, store db.Store, tokenMaker auth.Maker, config config.Config) {
	userHandler := handlers.NewUserHandler(store, tokenMaker, config)

	// Create a single user
	router.POST("/register", userHandler.CreateUser)
	// User login
	router.POST("/login", userHandler.Login)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
