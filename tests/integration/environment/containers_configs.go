package environment

import (
	"fmt"
	"path/filepath"
	"time"
)

const (
	NetworkName = "test_reports_service_network"

	AppName = "app"

	LocalStackName  = "localstack"
	LocalStackImage = "localstack/localstack"
	LocalStackPort  = "4566"

	AwsOpenSearchRegion    = "us-east-1"
	AwsOpenSearchDomain    = "reports-service-logs"
	AwsOpenSearchIndexName = "my-index"
	AwsOpenSearchEndpoint  = "reports-service-logs.us-east-1.es.localhost.localstack.cloud:4566"

	MigratorImage = "todo"        // todo
	MigratorName  = "es_migrator" // todo

	ElasticSearchName  = "elasticsearch"
	ElasticSearchImage = "docker.elastic.co/elasticsearch/elasticsearch:7.11.0"
	ElasticSearchPort  = "9200"

	KafkaImage              = "confluentinc/cp-kafka:7.2.0"
	KafkaDockerInternalPort = "29092" // https://www.confluent.io/blog/kafka-listeners-explained/ -> @see HOW TO: Connecting to Kafka on Docker
	KafkaHostExternalPort   = "9092"
	KafkaJMXClientPort      = "9999"
	KafkaContainerName      = "broker"

	ZookeeperImage         = "confluentinc/cp-zookeeper:7.2.0"
	ZooKeeperPort          = "2181"
	ZooKeeperContainerName = "zookeeper"
	ZooTickTime            = "2000"
)

var (
	pollInterval = time.Millisecond * 100
)

// doubledPort - helper func for prepare port mapping.
func doubledPort(port string) string {
	return fmt.Sprintf("%s:%[1]s", port) // output: "<port>:<port>"
}

// getAbsPath - get absolute path based on pwd and relative path.
func getAbsPath(relativePath string) string {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("could not resolve abs path for: %s. err is: %v", relativePath, err)
	}

	return path
}
