package transport

import (
	"errors"
	"log"
	"net/http"

	"image-reports/reporter/pkg/service"

	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
)

type serverConfiguration struct {
}

func NewServerConfiguration() server.ServerConfiguration[service.Service] {
	return &serverConfiguration{}
}

func (s *serverConfiguration) InitApiServer(router *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	return srv
}

func (s *serverConfiguration) InitUserService() service.Service {
	return service.NewService()
}

func (s *serverConfiguration) InitApiRoutes(svc service.Service) *gin.Engine {
	router := gin.Default()

	root := router.Group("/v1/reports")

	// Health check
	root.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
