package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"image-reports/reporter/pkg/endpoint"
	"image-reports/reporter/pkg/service"

	"image-reports/helpers/services/kafka"
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

	s.initKafkaListeners()

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
				log.Println(fmt.Errorf("could not read message on image processed: %w", err))
				continue
			}

			if err := endpoint.OnImageProcessedMessage(ctx, req); err != nil {
				log.Println(fmt.Errorf("could not handle message on image processed: %w", err))
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
				log.Println(fmt.Errorf("could not read message on image processed: %w", err))
				continue
			}

			if err := endpoint.OnImageStoredMessage(ctx, req); err != nil {
				log.Println(fmt.Errorf("could not handle message on image processed: %w", err))
			}
		}
	}()
}
