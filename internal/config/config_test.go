package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	const pathToConf = "../../config.conf"
	cfg, err := New(pathToConf + "123")

	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "could not parse config file."))
	} else {
		t.Fail()
	}
	assert.Nil(t, cfg)

	cfg, err = NewTestConfig(pathToConf)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.True(t, cfg.IsDevelopmentEnv())
	assert.False(t, cfg.IsProductionEnv())

	assert.Equal(t, "localhost:9092", cfg.GetString("kafka.addr"))
	assert.False(t, cfg.GetBoolOrDefaultValue("kafka.use_local_ca_cert", false))
	assert.Equal(t, "some", cfg.GetStringOrDefaultValue("kafka.use_local_ca_cert1", "some"))
	assert.EqualValues(t, "[stdout]", cfg.GetArray("logger.outputs").String())
	assert.Equal(t, 777, cfg.GetIntOrDefaultValue("kafka.timeout_ms1", 777))
	assert.Equal(t, true, cfg.GetBoolOrDefaultValue("kafka.timeout_ms1", true))
	assert.Equal(t, "development", cfg.GetCurrentEnvironment())

	key, err := cfg.GetJWTPrivateKey()
	assert.NotNil(t, key)
	assert.NoError(t, err)
}
