package transport

import (
	"errors"
	"fmt"
	"net/http"

	"image-reports/api-gateway/configs"
	processing_client "image-reports/api-gateway/helpers/http-client/processing"
	reporters_client "image-reports/api-gateway/helpers/http-client/reporter"
	storage_client "image-reports/api-gateway/helpers/http-client/storage"
	users_client "image-reports/api-gateway/helpers/http-client/users"
	"image-reports/api-gateway/pkg/endpoint"
	"image-reports/api-gateway/pkg/service"
	"image-reports/shared/models"

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
	return service.NewService(
		users_client.NewHttpClient(s.config.GlobalConfig),
		reporters_client.NewHttpClient(s.config.GlobalConfig),
		processing_client.NewHttpClient(s.config.GlobalConfig),
		storage_client.NewHttpClient(s.config.GlobalConfig),
	)
}

func (s *serverConfiguration) InitApiRoutes(svc service.Service) *gin.Engine {
	router := gin.Default()

	root := router.Group("/v1")
	root.Use(server.CORSMiddleware())

	// Health check
	root.GET("/ping/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authentication := root.Group("/auth")
	authentication.Use(server.CORSMiddleware())

	authentication.OPTIONS("/login")
	authentication.POST("/login", endpoint.Login(svc))

	resources := root.Group("/")
	resources.Use(
		server.CORSMiddleware(),
		auth.Authentication(),
		AddContextToken(),
		CheckUserValidity(svc),
	)

	reports := resources.Group("/reports")
	reports.Use(
		server.CORSMiddleware(),
		server.JSONMiddleware(),
	)

	reports.OPTIONS("/")
	reports.OPTIONS("/:id")
	reports.GET("/", endpoint.ListReports(svc))
	reports.GET("/:id", endpoint.GetReport(svc))
	reports.POST("/", endpoint.CreateReport(svc))
	reports.PATCH("/:id", auth.AllowOnlyRole(models.AdminRole), endpoint.ReportApproval(svc))

	storage := resources.Group("/storage")
	storage.Use(server.CORSMiddleware())

	storage.OPTIONS("/:id")
	storage.GET("/:id", endpoint.GetFile(svc))

	return router
}
