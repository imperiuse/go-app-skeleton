package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/gurkankaymak/hocon"
)

const AppName = "reports-service"

const (
	Development = "development"
	Production  = "production"
)

// Config - alias for clean deps.
type Config struct {
	*hocon.Config
}

// New - create new Config (hocon).
func New(configPath string) (*Config, error) {
	c, err := hocon.ParseResource(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file. err: %w", err)
	}

	return &Config{Config: c}, nil
}

// GetPortalJWTPublicKey - get portal public JWT key.
func (c *Config) GetPortalJWTPublicKey() ([]byte, error) {
	encodedPublicKey := c.GetString("portal.jwt_public_key")
	decoded, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		return []byte{}, fmt.Errorf("decode error: %w, \n input value was %v", err, encodedPublicKey)
	}

	return decoded, nil
}

// NewTestConfig - create new Config (hocon) ONLY FO TESTS NOT FOR PRODUCTION.
// //nolint: lll // this is for tests only.
func NewTestConfig(configPath string) (*Config, error) {
	setEnvs := [][2]string{
		{"KAFKA_ADDRESS", "localhost:9092"},
		{"KAFKA_CLIENT_PASSWORD", ""},
		{"KAFKA_TOPIC_REPORTS", "reports"},
		{"CURRENT_ENV", "development"},
		{"POSTGRES_HOST", "localhost"},
		{"POSTGRES_DB", "report_service"},
		{"POSTGRES_PORT", "5432"},
		{"POSTGRES_USER", "admin"},
		{"POSTGRES_PASSWORD", "password"},
		{"DISABLE_AUTH", "true"},
	}

	var err error
	for _, v := range setEnvs {
		if err = os.Setenv(v[0], v[1]); err != nil {
			return nil, err
		}
	}

	return New(configPath)
}

// GetString - like GetString redefine standard library get string. Unquote string if necessary.
func (c *Config) GetString(path string) string {
	s := c.Config.GetString(path)
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) && len(s) > 2 {
		return s[1 : len(s)-1]
	}

	return s
}

// GetPostgresDSN - get postgres dsn string.
func (c *Config) GetPostgresDSN() string {
	sslModeOption := ""

	postgresSSLMode := c.GetStringOrDefaultValue("postgres.ssl_mode", "disable")
	if postgresSSLMode != "" {
		sslModeOption = "sslmode=" + postgresSSLMode + " "
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %sTimeZone=Asia/Tokyo",
		c.GetString("postgres.host"), c.GetString("postgres.user"),
		c.GetString("postgres.password"), c.GetString("postgres.db"),
		c.GetInt("postgres.port"), sslModeOption,
	)
}

// todo re-write to generics when this proposal will be ready -> https://github.com/golang/go/issues/45380

// GetIntOrDefaultValue - get int if exist or return default value.
func (c *Config) GetIntOrDefaultValue(path string, defaultValue int) int {
	return getValueIfExist(c, path, c.GetInt, defaultValue)
}

// GetBoolOrDefaultValue - get bool if exist or return default value.
func (c *Config) GetBoolOrDefaultValue(path string, defaultValue bool) bool {
	return getValueIfExist(c, path, c.GetBoolean, defaultValue)
}

// GetStringOrDefaultValue - get string if exist or return default value.
func (c *Config) GetStringOrDefaultValue(path string, defaultValue string) string {
	return getValueIfExist(c, path, c.GetString, defaultValue)
}

// getValueIfExist - generic func for get T value from config struct.
func getValueIfExist[T any](c *Config, path string, f func(string) T, defaultValue T) T {
	parsedValue := c.Get(path)

	if parsedValue == nil {
		return defaultValue
	}

	return f(path)
}

// IsProductionEnv - is current env Production.
func (c *Config) IsProductionEnv() bool {
	return c.GetCurrentEnvironment() == Production
}

// IsDevelopmentEnv - is current env Development.
func (c *Config) IsDevelopmentEnv() bool {
	return c.GetCurrentEnvironment() == Development
}

// IsRunningInContainer - Check if app run in container.
func (c *Config) IsRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}
	return true
}

// GetCurrentEnvironment - get current env.
func (c *Config) GetCurrentEnvironment() string {
	return c.GetStringOrDefaultValue("environment", Production)
}

// GetJWTPrivateKey - get private JWT key.
func (c *Config) GetJWTPrivateKey() ([]byte, error) {
	encodedPrivateKey := c.GetString("jwt_private_key_base64")
	if encodedPrivateKey == "" {
		return []byte{}, errors.New("empty encodedPrivateKey. check config: `portal.jwt_private_key_base64`")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedPrivateKey)

	if err != nil {
		if c.IsDevelopmentEnv() {
			return []byte{}, fmt.Errorf("decode error: %w, \n input value was %v", err, encodedPrivateKey)
		}

		return []byte{}, fmt.Errorf("decode error: %w.   len: %d", err, len(encodedPrivateKey))
	}

	return decoded, nil
}

// ParseKafkaUserCredentials - return Username and Password for kafka user.
func (c *Config) ParseKafkaUserCredentials() (Username string, Password string, err error) {
	kafkaCredentials := struct {
		Address  string `json:"address"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err = jsoniter.Unmarshal([]byte(c.GetString("kafka.client_password")), &kafkaCredentials)
	if err != nil {
		return "", "", fmt.Errorf("could not Unmarshal kafkaCredentials from kafka.client_password")
	}

	return kafkaCredentials.Username, kafkaCredentials.Password, err
}
