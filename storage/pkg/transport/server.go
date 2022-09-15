package transport

import (
	"errors"
	"fmt"
	"net/http"

	"image-reports/storage/configs"
	"image-reports/storage/pkg/endpoint"
	"image-reports/storage/pkg/service"

	"image-reports/helpers/services/auth"
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
	addr := fmt.Sprintf("%s:%d", s.config.Services.Storage.Host, s.config.Services.Storage.Port)
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

	storage := root.Group("/storage")

	storage.Use(auth.Authentication())

	storage.GET("/:id", endpoint.GetImage(s.config.Path))
	storage.POST("/",
		server.JSONMiddleware(),
		endpoint.SaveImage(s.config.Path, s.config.Image.MaxSize, s.config.Image.Extensions),
	)

	return router
}
