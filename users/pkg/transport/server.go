package transport

import (
	"errors"
	"fmt"
	"net/http"

	"image-reports/shared/models"
	"image-reports/users/configs"
	"image-reports/users/pkg/endpoint"
	"image-reports/users/pkg/service"

	"image-reports/helpers/services/auth"
	log "image-reports/helpers/services/logger"
	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type serverConfiguration struct {
	db     *gorm.DB
	config *configs.AppConfig
}

func NewServerConfiguration(db *gorm.DB, config *configs.AppConfig) server.ServerConfiguration[service.Service] {
	return &serverConfiguration{
		db:     db,
		config: config,
	}
}

func (s *serverConfiguration) InitApiServer(router *gin.Engine) *http.Server {
	addr := fmt.Sprintf("%s:%d", s.config.Services.Users.Host, s.config.Services.Users.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Infof("Running server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()

	return srv
}

func (s *serverConfiguration) InitUserService() service.Service {
	return service.NewService(s.db)
}

func (s *serverConfiguration) InitApiRoutes(svc service.Service) *gin.Engine {
	router := gin.Default()

	root := router.Group("/v1")

	// Health check
	root.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authentication := root.Group("/auth")
	authentication.Use(server.JSONMiddleware())

	authentication.POST("/", endpoint.CheckCredentials(svc))
	authentication.GET("/:id", auth.Authentication(), endpoint.CheckUserId(svc))

	users := root.Group("/users")
	users.Use(
		server.JSONMiddleware(),
		auth.Authentication(),
		auth.AllowOnlyRole(models.AdminRole),
	)

	users.GET("/:id", endpoint.GetUser(svc))
	users.POST("/search", endpoint.SearchUsers(svc))

	return router
}
