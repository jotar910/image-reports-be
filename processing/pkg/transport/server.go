package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"image-reports/processing/configs"
	"image-reports/processing/pkg/endpoint"
	"image-reports/processing/pkg/service"

	"image-reports/helpers/services/auth"
	"image-reports/helpers/services/kafka"
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
	addr := fmt.Sprintf("%s:%d", s.config.Services.Processing.Host, s.config.Services.Processing.Port)
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

	evaluation := root.Group("/evaluations")

	evaluation.Use(
		server.JSONMiddleware(),
		auth.Authentication(),
	)

	evaluation.GET("/:id", endpoint.GetEvaluation(svc))
	evaluation.POST("/", endpoint.ProcessImage(svc, s.config.Image.MaxSize, s.config.Image.Extensions))

	s.initKafkaListeners(svc)

	return router
}

func (s *serverConfiguration) initKafkaListeners(svc service.Service) {
	go func() {
		r := kafka.Reader(kafka.TopicImageProcessed, kafka.TopicImageProcessed)
		req := kafka.NewEmptyImageProcessedMessage()
		for {
			err := r.Read(context.Background(), req)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Errorf("could not read message on image processed: %w", err)
				continue
			}
			log.Debugf("received message about image processed: %+v", req)
			if err := endpoint.AddEvaluation(svc, req); err != nil {
				log.Errorf("could not add evaluation on image processed: %w", err)
				continue
			}
		}
	}()
}
