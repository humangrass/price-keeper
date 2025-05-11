package xhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RouteRegister interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Server struct {
	*http.Server
	logger *zap.SugaredLogger
	Mux    *http.ServeMux
}

func NewServer(opt Opt, logger *zap.SugaredLogger) (*Server, error) {
	err := opt.Validate()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	if opt.EnableHealthMethod {
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte("OK"))
			if err != nil {
				logger.Errorf("HTTP server error: %v", err)
			}
		})
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Handler:      mux,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
	}

	return &Server{
		Server: httpServer,
		logger: logger,
		Mux:    mux,
	}, nil
}

func (s *Server) Start() {
	go func() {
		s.logger.Infof("Starting HTTP server on %s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("HTTP server error: %v", err)
		}
	}()
}

func (s *Server) Drop() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info("Shutting down HTTP server...")
	if err := s.Shutdown(shutdownCtx); err != nil {
		s.logger.Errorf("HTTP server shutdown error: %v", err)
		return err
	}
	s.logger.Info("HTTP server stopped gracefully")
	return nil
}

func (s *Server) DropMsg() string {
	return "graceful shutdown of HTTP server"
}

func (s *Server) RegisterRoutes(registrator RouteRegister) {
	registrator.RegisterRoutes(s.Mux)
}
