package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	gormLogger "gorm.io/gorm/logger"

	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"

	_ "github.com/imperiuse/go-app-skeleton/docs"
	"github.com/imperiuse/go-app-skeleton/internal/config"
	"github.com/imperiuse/go-app-skeleton/internal/database"
	"github.com/imperiuse/go-app-skeleton/internal/database/migration"
	"github.com/imperiuse/go-app-skeleton/internal/database/tables"
	"github.com/imperiuse/go-app-skeleton/internal/logger"
	"github.com/imperiuse/go-app-skeleton/internal/logger/field"
	"github.com/imperiuse/go-app-skeleton/internal/servers/api"
	"github.com/imperiuse/go-app-skeleton/internal/servers/metrics"
	"github.com/imperiuse/go-app-skeleton/internal/servers/pprof"

	// Automatically set GOMAXPROCS to match Linux container CPU quota.
	_ "go.uber.org/automaxprocs"

	_ "net/http/pprof" //nolint
)

// Start and Finished time for App. Related for DI.
const (
	appStartTimeout = 60 * time.Second
	appStopTimeout  = 60 * time.Second
)

// Program flags.
// This flag only for dev env and developer's purposes (must not be used in prod).
var (
	configPath                   = flag.String("config", "config.conf", "path to config file (with HOCON format)")
	notRunMetricsAndPprofServers = flag.Bool("notRunMetricsAndPprofServers", false,
		"No run metrics and pprof servers (this might be helpfully to resolve port conflict when debug complex envs with id-prov)")
)

type application struct {
	version    string
	configPath string

	startTimeout time.Duration
	stopTimeout  time.Duration
}

// @title Reports service Swagger HTTP API
// @version 1.0.0
// @description This is Swagger docs for HTTP REST API reports service

// @contact.name API Support
// @contact.email arseny.sazanov@gmail.com

// @license.name None
// @license.url None

// @host localhost:8080
// @BasePath /
func main() {
	flag.Parse()

	app := &application{
		version:      os.Getenv("APP_VERSION"),
		configPath:   *configPath,
		startTimeout: appStartTimeout,
		stopTimeout:  appStopTimeout,
	}

	app.run()
}

// //nolint:funlen // this is ok, dirty func for DI purposes.
func (a *application) run() {
	fxApp := fx.New(
		fx.Provide(
			func() (context.Context, context.CancelFunc, *errgroup.Group) {
				ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
				g, gCtx := errgroup.WithContext(ctx)

				return gCtx, cancel, g
			},
			func() (*config.Config, error) {
				return config.New(a.configPath)
			},
			func(cfg *config.Config) (*logger.Logger, error) {
				return logger.New(logger.Config{
					Level:    cfg.GetString("logger.level"),
					Encoding: cfg.GetString("logger.encoding"),
					Color:    cfg.GetBoolean("logger.color"),
					Outputs:  cfg.GetStringSlice("logger.outputs"),
					Tags:     cfg.GetStringSlice("logger.tags"),
				}, cfg.GetCurrentEnvironment(), config.AppName, a.version)
			},
			func(cfg *config.Config, log *logger.Logger) gormLogger.Interface {
				var lvl = gormLogger.Warn
				if cfg.IsDevelopmentEnv() {
					lvl = gormLogger.Warn
				}

				return logger.NewGormLogger(log, gormLogger.Config{
					SlowThreshold:             30 * time.Second,
					Colorful:                  false,
					IgnoreRecordNotFoundError: true,                  // if true - not found - is not error.
					ParameterizedQueries:      cfg.IsProductionEnv(), // when True - Don't include params in the SQL log
					LogLevel:                  lvl,
				})
			},
			api.NewEngine,
			func(cfg *config.Config, log *logger.Logger) *pprof.Server {
				return pprof.New(pprof.Config{
					Name:    "pprof",
					Address: cfg.GetString("servers.pprof.addr"),
				}, log)
			},
			func(cfg *config.Config, log *logger.Logger) *metrics.Server {
				return metrics.New(metrics.Config{
					Name:    "metrics",
					Address: cfg.GetString("servers.metrics.addr"),
				}, log)
			},
			func(cfg *config.Config, e *api.Engine, log *logger.Logger) *api.Server {
				return api.NewServer(api.Config{
					IsDevEnv:       cfg.IsDevelopmentEnv(),
					ServiceName:    config.AppName,
					Addr:           cfg.GetString("servers.api.addr"),
					DisableAuth:    cfg.GetBoolOrDefaultValue("servers.api.disable_auth", false),
					AllowOrigin:    cfg.GetString("servers.api.allow_origin"),
					EnableStatsViz: cfg.GetBoolean("servers.api.enable_statsviz"),
					WriteTimeout:   cfg.GetDuration("servers.api.write_timeout"),
					ReadTimeout:    cfg.GetDuration("servers.api.read_timeout"),
				},
					e, log,
				)
			},
			func(
				cfg *config.Config,
				log *logger.Logger,
				gLogger gormLogger.Interface,
				configuration *config.Config,
				shutdowner fx.Shutdowner,
			) (*database.DB, error) {
				return database.New(configuration.GetPostgresDSN(), false, gLogger, shutdowner)
			},
		),

		fx.Invoke(func(log *logger.Logger, db *database.DB, configuration *config.Config) error {
			// *repl.ReplService - is needed because we need run migrations after repl service migration.
			var err error
			if err = migration.ApplyMigrations(db, tables.AllDTOs[:]...); err != nil {
				return fmt.Errorf("migration has not applied: %w", err)
			}

			return nil
		},
			a.start),
		fx.StartTimeout(a.startTimeout),
		fx.StopTimeout(a.startTimeout),
	)

	fxApp.Run()
}

func (a *application) start(
	lc fx.Lifecycle,
	gCtx context.Context,
	gCancel context.CancelFunc,
	cfg *config.Config,
	errGroup *errgroup.Group,
	log *logger.Logger,
	db *database.DB,
	mServer *metrics.Server,
	pServer *pprof.Server,
	apiServer *api.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("App starting successfully! Ready for hard work!",
				field.String("version", a.version), field.String("appName", config.AppName))

			if !*notRunMetricsAndPprofServers {
				mServer.Run()
				pServer.Run()
			}

			if cfg.IsDevelopmentEnv() {
				log.Sugar().Info("Kibana (Elastics web UI viewer) available here -> http://localhost:5601/app/home#/ " +
					"More details you can find in docker compose file -> `docker/docker-compose-dev-local.yml` section `kibana`")
			}

			apiServer.Run(errGroup, gCtx, appStopTimeout)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			gCancel()

			err := apiServer.Stop(ctx)

			log.Info("App Finishing! ...",
				field.String("version", a.version), field.String("appName", config.AppName))

			return err
		},
	})
}
