package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/softstone1/fl/config"
)

// Server encapsulates the HTTP server and its dependencies
type Server struct {
	srv *http.Server
}

// NewServer initializes a new Server with the provided configuration
func NewServer(cfg *config.Config, router http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
			Handler:      router,
			ReadTimeout:  cfg.ServerReadTimeout,
			WriteTimeout: cfg.ServerWriteTimeout,
		},
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	// This blocks until an error occurs, or srv.Shutdown() is called
	return s.srv.ListenAndServe()
}
