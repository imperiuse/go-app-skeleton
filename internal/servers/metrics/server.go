package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/imperiuse/go-app-skeleton/internal/logger"
)

type (
	// Config - Config.
	Config struct {
		Name    string
		Address string
	}

	// Server - Server.
	Server struct {
		config Config
		server *http.Server
		logger *logger.Logger
	}
)

// New create server.
func New(config Config, log *logger.Logger) *Server {
	return &Server{
		config: config,
		logger: log,
		server: &http.Server{
			Addr:           config.Address,
			Handler:        promhttp.Handler(),
			ReadTimeout:    time.Minute,
			WriteTimeout:   time.Minute,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Run - run server.
func (s *Server) Run() {
	go func() {
		s.logger.Sugar().Infof("starting metrics server at->  http://localhost%s", s.config.Address)

		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Fatal("error while starting metrics server. Is port free?")
		}
	}()
}
