package pprof

import (
	"net/http"
	"runtime"
	"time"

	"github.com/imperiuse/go-app-skeleton/internal/logger"
)

const (
	mutexProfileFraction = 5
	blockProfileRate     = 100
)

type (
	// Config - config pprof server.
	Config struct {
		Name    string
		Address string
	}

	// Server - pprof server.
	Server struct {
		config Config
		server *http.Server
		logger *logger.Logger
	}
)

// New - new pprof server.
func New(config Config, log *logger.Logger) *Server {
	return &Server{
		config: config,
		logger: log,
		server: &http.Server{
			Addr:           config.Address,
			ReadTimeout:    time.Minute,
			WriteTimeout:   time.Minute,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Run - start pprof server.
func (s *Server) Run() {
	go func() {
		s.logger.Sugar().Infof("starting pprof server at->  http://localhost%s/debug/pprof/", s.config.Address)

		runtime.SetMutexProfileFraction(mutexProfileFraction)
		runtime.SetBlockProfileRate(blockProfileRate)

		if err := s.server.ListenAndServe(); err != nil {
			s.logger.Fatal("error while starting pprof server. Is port free?")
		}
	}()
}
