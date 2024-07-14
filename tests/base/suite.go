package base

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/imperiuse/go-app-skeleton/internal/config"
	"github.com/imperiuse/go-app-skeleton/internal/database"
	"github.com/imperiuse/go-app-skeleton/internal/logger"
	"github.com/imperiuse/go-app-skeleton/internal/servers/api"
)

type (
	SuiteForUnit struct {
		suite.Suite

		Ctx    context.Context
		Cancel context.CancelFunc

		Configuration *config.Config

		Engine                  *api.Engine
		DefaultTenantID         uuid.UUID
		DefaultRequestingUserID uuid.UUID
		Server                  *api.Server

		DB *database.DB
	}

	SuiteForIntegration struct {
		SuiteForUnit
	}
)

type (
	TestCaseHTTP struct {
		RequestProvider    func() (*http.Request, error)
		ExpectedStatusCode int
		ExpectedBody       string
	}
)

func (s *SuiteForUnit) SetupConfig() {
	s.NoError(os.Setenv("KAFKA_ADDRESS", "broker:29092"))
	s.NoError(os.Setenv("KAFKA_CLIENT_PASSWORD", ""))
	s.NoError(os.Setenv("CURRENT_ENV", "development"))
	s.NoError(os.Setenv("AWS_OPEN_SEARCH_SERVICE_HOST", "http://reports-service-logs.us-east-1.es.localhost.localstack.cloud:4566"))
	s.NoError(os.Setenv("AWS_OPEN_SEARCH_SERVICE_INDEX", "my-index"))

	var err error
	s.Configuration, err = config.NewTestConfig("../../../config.conf")
	s.Nil(err, "error must be nil for hocon.ParseResource")
	s.NotNil(s.Configuration, "obj cfg must be not nil")
}

func (s *SuiteForUnit) Setup() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	s.SetupConfig()

	var err error
	s.DB, err = database.NewWithoutFX(s.Configuration.GetPostgresDSN(), true,
		logger.NewGormLogger(zap.NewNop(), logger.GormLoggerConfig{
			SlowThreshold:             time.Second * 30,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  gormLogger.Warn,
		}))
	s.Nil(err, "err create DB")
	s.NotNil(s.DB, "DB must be not nil")

	s.Engine = api.NewEngine()
	s.Server = api.NewServer(api.Config{
		IsDevEnv:       s.Configuration.IsDevelopmentEnv(),
		ServiceName:    "test-api-healthHandlersTestSuite",
		Addr:           s.Configuration.GetString("server.addr"),
		DisableAuth:    s.Configuration.GetBoolean("server.disable_auth") && s.Configuration.IsDevelopmentEnv(),
		EnableStatsViz: s.Configuration.GetBoolean("server.enable_statsviz"),
		WriteTimeout:   s.Configuration.GetDuration("server.write_timeout"),
		ReadTimeout:    s.Configuration.GetDuration("server.read_timeout"),
	},
		s.Engine,
		zap.NewNop(),
	)
}

func (s *SuiteForUnit) AssertResultOfHTTPRequest(testCase TestCaseHTTP) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := testCase.RequestProvider()
	if err != nil {
		s.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	s.Engine.ServeHTTP(rr, req)

	s.T().Logf("Make requests: %+v", req)

	// Check the status code is what we expect.
	s.Equal(testCase.ExpectedStatusCode, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	s.Equal(testCase.ExpectedBody, rr.Body.String(), "handler returned unexpected body")
}

func (s *SuiteForIntegration) Setup() {
	s.SetupConfig()
}
