package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"image-reports/reporter/configs"
	"image-reports/reporter/pkg/endpoint"
	"image-reports/reporter/pkg/service"

	"image-reports/helpers/services/kafka"
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
	addr := fmt.Sprintf("%s:%d", s.config.Services.Reporter.Host, s.config.Services.Reporter.Port)
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
				log.Errorf("could not read message on image processed: %v", err)
				continue
			}

			if err := endpoint.OnImageProcessedMessage(ctx, req); err != nil {
				log.Errorf("could not handle message on image processed: %v", err)
			}
		}
	}()

	go func() {
		r := kafka.Reader(kafka.TopicImageStored, kafka.TopicImageStoredGroup)
		req := kafka.NewEmptyImageStoredMessage()
		for {
			ctx := context.Background()
			err := r.Read(ctx, req)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Errorf("could not read message on image processed: %v", err)
				continue
			}

			if err := endpoint.OnImageStoredMessage(ctx, req); err != nil {
				log.Errorf("could not handle message on image processed: %v", err)
			}
		}
	}()
}
