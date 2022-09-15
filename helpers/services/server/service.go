package server

import (
	"context"
	"image-reports/helpers/services/kafka"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "image-reports/helpers/services/logger"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run()
	Shutdown(ctx context.Context)
}

type ServerConfiguration[TService any] interface {
	InitApiServer(router *gin.Engine) *http.Server
	InitUserService() TService
	InitApiRoutes(svc TService) *gin.Engine
}

type serverTemplate[TService any] struct {
	ServerConfiguration[TService]
	srv *http.Server
}

func NewServer[TService any](config ServerConfiguration[TService]) Server {
	return &serverTemplate[TService]{config, nil}
}

func (s *serverTemplate[TService]) Run() {
	svc := s.InitUserService()
	router := s.InitApiRoutes(svc)
	s.srv = s.InitApiServer(router)

	// Wait for a signal that ends the execution.
	signc := make(chan os.Signal, 1)
	signal.Notify(signc, os.Interrupt)

	sign := <-signc
	log.Infof("Receive terminate, grateful shutdown", sign)

	// Set a timeout to end the server gracefully.
	tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// End gracefully.
	s.Shutdown(tc)
	os.Exit(0)
}

func (s *serverTemplate[TService]) Shutdown(ctx context.Context) {
	kafka.Shutdown()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}
}
