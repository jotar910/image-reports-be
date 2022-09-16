package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"image-reports/realtime/configs"
	"image-reports/realtime/pkg/endpoint"
	"image-reports/realtime/pkg/service"

	"image-reports/helpers/services/kafka"
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
	addr := fmt.Sprintf("%s:%d", s.config.Services.Realtime.Host, s.config.Services.Realtime.Port)
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

	s.initKafkaListeners()

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

	return router
}

func (s *serverConfiguration) initKafkaListeners() {
	go func() {
		r := kafka.Reader(kafka.TopicImageProcessed, kafka.TopicImageProcessedGroup)
		req := kafka.NewEmptyImageProcessedMessage()
		for {
			ctx := context.Background()
			err := r.Read(ctx, req)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Errorf("could not read message on image processed: %w", err)
				continue
			}

			if err := endpoint.OnImageProcessedMessage(ctx, req); err != nil {
				log.Errorf("could not handle message on image processed: %w", err)
			}
		}
	}()
}
