package transport

import (
	"errors"
	"fmt"
	"net/http"

	"image-reports/api-gateway/configs"
	"image-reports/api-gateway/pkg/endpoint"
	"image-reports/api-gateway/pkg/service"

	users_client "image-reports/helpers/services/http-client/users"
	log "image-reports/helpers/services/logger"
	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
)

type serverConfiguration struct {
	config *configs.AppConfig
}

func NewServerConfiguration(config *configs.AppConfig) server.ServerConfiguration[service.Service] {
	return &serverConfiguration{
		config: config,
	}
}

func (s *serverConfiguration) InitApiServer(router *gin.Engine) *http.Server {
	addr := fmt.Sprintf("%s:%d", s.config.Services.ApiGateway.Host, s.config.Services.ApiGateway.Port)
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
	return service.NewService(users_client.NewHttpClient(s.config.GlobalConfig))
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

	auth := root.Group("/auth")
	auth.POST("/login", endpoint.Login(svc))

	return router
}
