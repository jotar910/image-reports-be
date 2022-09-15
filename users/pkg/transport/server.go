package transport

import (
	"errors"
	"fmt"
	"net/http"

	"image-reports/users/configs"
	"image-reports/users/pkg/endpoint"
	"image-reports/users/pkg/service"

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

	root.POST("/auth", endpoint.CheckCredentials(svc))

	return router
}
