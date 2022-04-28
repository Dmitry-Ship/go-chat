package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	httpServer *http.Server
}

func NewGracefulServer() *server {
	port := os.Getenv("PORT")

	httpServer := &http.Server{
		Addr: ":" + port,
	}
	return &server{
		httpServer: httpServer,
	}
}

func (s *server) Run() {
	defer s.shutdownGracefully()

	go func() {
		log.Println("Listening on port " + s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	s.waitForStopSignal()
}

func (s *server) shutdownGracefully() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
}

func (s *server) waitForStopSignal() {
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
