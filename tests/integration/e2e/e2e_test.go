package e2e

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/imperiuse/go-app-skeleton/tests/base"
	"github.com/imperiuse/go-app-skeleton/tests/integration/environment"
)

type End2EndTestSuite struct {
	base.SuiteForIntegration

	environment.ContainersEnvironment
}

// TestEnd2EndTestSuite - Root test for test suite End2EndTestSuite.
func TestEnd2EndTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(End2EndTestSuite))
}

func (s *End2EndTestSuite) SetupSuite() {
	const magicTime = 3 * time.Second

	s.T().Log("> From SetupSuite")

	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	composeFilePaths := []string{
		"../../../../docker/docker-compose-e2e.yml",
	}
	identifier := strings.ToLower(uuid.New().String())

	s.StartDockerComposeEnvironment(s.T(), composeFilePaths, identifier)

	s.T().Log("Start setup app services")

	s.SuiteForIntegration.Setup()

	<-time.After(magicTime)
}

func (s *End2EndTestSuite) SetupTest() {
	s.T().Log(">> From SetupTest")
}

func (s *End2EndTestSuite) BeforeTest(_, _ string) {
	s.T().Log(">>> From BeforeTest")
}

func (s *End2EndTestSuite) AfterTest(_, _ string) {
	s.T().Log(">>> From AfterTest")
}

func (s *End2EndTestSuite) TearDownTest() {
	s.T().Log(">> From TearDownTest")
}

func (s *End2EndTestSuite) TearDownSuite() {
	s.T().Log("> From TearDownSuite")

	s.FinishedDockerComposeEnvironment(s.T())

	s.Cancel()
}

func (s *End2EndTestSuite) TestBaseCheckAppStartedSuccessfully() {
	s.T().Log(">>>> From TestBaseCheckAppStartedSuccessfully")

	res, err := http.Get("http://localhost:8080/health")
	require.Nilf(s.T(), err, "err must be nil. but got error: %v", err)
	require.Equal(s.T(), http.StatusOK, res.StatusCode, "must be equal")
	body, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err, "err must be nil")
	require.Equal(s.T(), `{"alive":true}`, string(body), "must be equal")
	defer res.Body.Close()

	res, err = http.Get("http://localhost:8080/ready")
	require.Nilf(s.T(), err, "err must be nil. but got error: %v", err)
	require.Equal(s.T(), http.StatusOK, res.StatusCode, "must be equal")
	body, err = io.ReadAll(res.Body)
	require.NoError(s.T(), err, "err must be nil")
	require.Equal(s.T(), `{"ready":true}`, string(body), "must be equal")
	defer res.Body.Close()
}
