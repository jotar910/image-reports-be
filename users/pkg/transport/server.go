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
)

type serverConfiguration struct {
	server.ServerConfigurationHooks[service.Service]
	config *configs.AppConfig
}

func NewServerConfiguration() server.ServerConfiguration[service.Service] {
	return &serverConfiguration{}
}

func (s *serverConfiguration) BeforeInit() {
	if _, err := configs.Initialize("users"); err != nil {
		log.Fatalf("config: %s", err)
	}
	s.config = configs.Get()
	log.Infof("Starting with config: %+v", s.config)

	gin.SetMode(s.config.Gin.Mode)
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
	return service.NewService()
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
