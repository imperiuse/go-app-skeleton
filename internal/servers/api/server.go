package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// Swagger embed files.
	filesSwagger "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/net/webdav"
	"golang.org/x/sync/errgroup"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"

	"github.com/imperiuse/go-app-skeleton/internal/logger"
	"github.com/imperiuse/go-app-skeleton/internal/logger/field"
	"github.com/imperiuse/go-app-skeleton/internal/metrics"
	mw "github.com/imperiuse/go-app-skeleton/internal/servers/api/middleware"
)

type (
	// Config - config for http server.
	Config struct {
		IsDevEnv       bool
		ServiceName    string
		Addr           string
		DisableAuth    bool
		EnableStatsViz bool
		AllowOrigin    string
		WriteTimeout   time.Duration
		ReadTimeout    time.Duration
	}

	// Server - http API server structure.
	Server struct {
		server    *http.Server
		ginEngine *gin.Engine
		log       *logger.Logger
	}
)

// NewServer - constructor http API Server.
func NewServer(cfg Config, e *Engine, log *logger.Logger) *Server {
	if !cfg.IsDevEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		server: &http.Server{
			Addr:           cfg.Addr,
			Handler:        e,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		ginEngine: e,
		log:       log,
	}

	// Setup main middleware
	log.Info("Starting create middleware and routes for gin server")

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	e.Use(mw.Ginzap(log, time.RFC3339, true))

	// Add Prometheus metrics for endpoint
	e.Use(metrics.PrometheusMiddleware())

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	e.Use(mw.RecoveryWithZap(log, true))

	// https://stackoverflow.com/questions/29418478/go-gin-framework-cors
	e.Use(mw.CORSMiddleware(cfg.AllowOrigin))

	e.Use(mw.UUIDMiddleware())

	e.Use(otelgin.Middleware(cfg.ServiceName))

	// add statsviz (viewer of pprof) @see more here -> https://github.com/arl/statsviz
	if cfg.IsDevEnv {
		// Create statsviz server.
		srv, _ := statsviz.NewServer()

		ws := srv.Ws()
		index := srv.Index()

		e.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
			if context.Param("filepath") == "/ws" {
				ws(context.Writer, context.Request)

				return
			}
			index(context.Writer, context.Request)
		})

		log.Debug(fmt.Sprintf("start statsviz at -> http://localhost%v/%v", cfg.Addr, "debug/statsviz"))
	}

	// Swagger Docs.
	e.GET("/swagger/*any", DisablingWrapHandler(filesSwagger.Handler, !cfg.IsDevEnv))

	e.GET("/health", healthHandler)
	e.GET("/ready", readyHandler)

	apiV1 := e.Group("/api/v1/")

	apiV1 = apiV1

	return s
}

// Run - start http API server.
func (s *Server) Run(g *errgroup.Group, gCtx context.Context, shutdownTimeout time.Duration) {
	s.log.Info("Starting http api server", field.String("addr", s.server.Addr))

	g.Go(func() error {
		return s.server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		s.log.Info("gCtx.Done. Shutdown http api server.")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return s.Stop(shutdownCtx)
	})
}

// Stop - stop server.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// DisablingWrapHandler - Disable documentation for production environment.
func DisablingWrapHandler(h *webdav.Handler, isDisabled bool) gin.HandlerFunc {
	if isDisabled {
		return func(c *gin.Context) {
			c.String(http.StatusNotFound, "")
		}
	}

	return ginSwagger.WrapHandler(h)
}

// Health godoc
// @Summary Health check
// @Description health check
// @Id Health
// @Tags Server Base
// @Accept  json
// @Produce  json
// @Success 200
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}

// Health godoc
// @Summary Ready check
// @Description ready check
// @Id Ready
// @Tags Server Base
// @Accept  json
// @Produce  json
// @Success 200
// @Router /ready [get]
func readyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}
