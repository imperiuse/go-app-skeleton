package environment

import (
	"context"
	"testing"
	"time"

	"github.com/imperiuse/go-app-skeleton/tests/integration/testcontainer"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type ContainersEnvironment struct {
	// First - native go-testcontainers way for creating docker containers.
	dockerNetwork     *testcontainer.DockerNetwork
	postgresContainer testcontainers.Container

	// Second - docker-compose way + go-testcontainers for create docker container env.
	// compose compose.DockerCompose
}

// StartPureDockerEnvironment - create and start docker containers env with first way.
func (c *ContainersEnvironment) StartPureDockerEnvironment(t *testing.T, ctx context.Context) {
	t.Log("Start Test Pure Docker based Environment")

	t.Log("Create docker network")
	dn, err := testcontainer.NewDockerNetwork(ctx, NetworkName)
	require.Nil(t, err, "error must be nil for NewDockerNetwork")
	require.NotNil(t, dn, "docker network must be not nil")
	c.dockerNetwork = dn.(*testcontainers.DockerNetwork)

	t.Log("Create service deps")

	t.Log("Start deps services containers")

	require.Nil(t, c.postgresContainer.Start(ctx), "postgres must start without errors")

	const magicTime = time.Second * 3
	time.Sleep(magicTime)
}

// FinishedPureDockerEnvironment - finished containers (env) which we created by first way.
func (c *ContainersEnvironment) FinishedPureDockerEnvironment(t *testing.T, ctx context.Context) {
	t.Log("Finished Test Pure Docker Environment from files")
	require.Nil(t, testcontainer.TerminateIfNotNil(ctx, c.postgresContainer), "must not get an error while terminate postgres cluster")
	require.Nil(t, c.dockerNetwork.Remove(ctx), "must not get an error while remove docker network")
}

// StartDockerComposeEnvironment - create and start docker containers env with second way.
func (c *ContainersEnvironment) StartDockerComposeEnvironment(
	t *testing.T,
	composeFilePaths []string,
	identifier string,
) {
	//t.Logf("Start Test Dockercompose based Environment from files: %+v", composeFilePaths)
	//c.compose = compose.NewLocalDockerCompose(composeFilePaths, identifier).
	//	WaitForService(ZooKeeperContainerName, wait.ForLog("binding to port 0.0.0.0/0.0.0.0:"+ZooKeeperPort)).
	//	WaitForService(KafkaContainerName, wait.ForLog("[KafkaServer id=1] started"))
	//
	//if len(composeFilePaths) > 1 { // this is little tricky hack here. :)
	//	// if we have one docker-compose file for app container, that add wait strategy.
	//	c.compose = c.compose.WaitForService(AppName, wait.ForLog("App starting successfully! Ready for hard work!"))
	//}
	//
	//require.Nil(t, c.compose.WithCommand([]string{"up", "--force-recreate", "-d"}).Invoke().Error)
}

// FinishedDockerComposeEnvironment - finished containers (env) which we created by second way.
func (c *ContainersEnvironment) FinishedDockerComposeEnvironment(t *testing.T) {
	t.Log("Finished Test Docker compose based Environment from files")
	//require.Nil(t, c.compose.Down().Error, "docker compose must down without errors")
}
