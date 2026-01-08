package routes

import (
	"github.com/SecureParadise/go_attendence/internal/api/middleware"
	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/SecureParadise/go_attendence/internal/config"
	"github.com/SecureParadise/go_attendence/internal/db"
	"github.com/SecureParadise/go_attendence/internal/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	store      db.Store
	tokenMaker auth.Maker
	router     *gin.Engine
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()

	// Middlewares
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(util.GinLogger())
	router.Use(gin.Recovery())

	// Setup routes
	SetupUnProtectedRoutes(router, server.store, server.tokenMaker, server.config)
	SetupProtectedRoutes(router, server.store, server.tokenMaker, server.config)

	server.router = router
}

func (server *Server) GetRouter() *gin.Engine {
	return server.router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
