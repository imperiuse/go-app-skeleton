package logger

import (
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/imperiuse/go-app-skeleton/internal/logger/field"
)

func Test_NewLogger(t *testing.T) {
	log, err := New(Config{}, "", "", "")
	assert.NotNil(t, err)
	assert.Nil(t, log)

	log, err = New(Config{
		Level:    "info",
		Encoding: "console",
		Color:    true,
		Outputs:  []string{"stdout"},
	}, "prod", "core", "v1.2.3")
	assert.Nil(t, err)
	assert.NotNil(t, log)

	log, err = New(Config{
		Level:    "debug",
		Encoding: "json",
		Color:    false,
		Outputs:  []string{"stdout"},
	}, "test", "core", "dev")
	assert.NotNil(t, log)
	assert.Nil(t, err)

	err = nil
	LogIfError(log, "message1", err, field.Any("testField", 123))
	err = fmt.Errorf("err")
	LogIfError(log, "message2", err, field.Any("testField", 123))
	// Output:
	// {"level":"ERROR","ts":"2021-01-31T23:15:26.052+0300",
	// "caller":"logger/logger.go:71","msg":"msg2","env":"test","version":"dev","services":"core","testField":123,"error":"err"}
}

func Test_Loggers_Fields(t *testing.T) {
	const f = "field"
	for _, v := range []any{"", 1, "1234", nil, struct {
		A int
	}{A: 1234}} {
		assert.Equal(t, zap.Any(f, v), field.Any(f, v))
	}
	assert.Equal(t, zap.Bool(f, true), field.Bool(f, true))
	assert.Equal(t, zap.Bool(f, false), field.Bool(f, false))
	assert.Equal(t, zap.Int(f, 123), field.Int(f, 123))
	assert.Equal(t, zap.Int64(f, int64(123)), field.Int64(f, 123))
	assert.Equal(t, zap.Uint64(f, 123), field.Uint64(f, uint64(123)))
	assert.Equal(t, zap.Int64("id", 123), field.ID(123))
	assert.Equal(t, zap.Int("port", 123), field.Port(123))
	assert.Equal(t, zap.Int("id", 123), field.ID(123))
	assert.Equal(t, zap.String("id", "123"), field.ID("123"))
	assert.Equal(t, zap.Int64("id", 123), field.ID(int64(123)))
	assert.Equal(t, zap.Int32("id", 123), field.ID(int32(123)))
	assert.Equal(t, zap.String("id", "id"), field.ID("id"))
	assert.Equal(t, zap.String(f, "123"), field.String(f, "123"))
	assert.Equal(t, zap.String("controller", "general"), field.String("controller", "general"))
	assert.Equal(t, zap.String("repo", "123"), field.Repo("123"))
	assert.Equal(t, zap.String("topic", "123"), field.Topic("123"))
	assert.Equal(t, zap.String("controller", "ccc"), field.Controller("ccc"))
	assert.Equal(t, zap.String("handler", "hhh"), field.Handler("hhh"))
	assert.Equal(t, zap.String("table", "123"), field.Table("123"))
	assert.Equal(t, zap.String("env", "123"), field.Env("123"))
	assert.Equal(t, zap.String("version", "123"), field.Version("123"))
	assert.Equal(t, zap.String("service", "123"), field.Service("123"))
	assert.Equal(t, zap.String("service", "123"), field.Service("123"))
	assert.Equal(t, zap.String("host", "123"), field.Host("123"))
	assert.Equal(t, zap.String("path", "123"), field.Path("123"))
	assert.Equal(t, zap.String("ip", "123"), field.IP("123"))
	assert.Equal(t, zap.String("ip_forwarded_for", "123"), field.IPForwardedFor("123"))
	assert.Equal(t, zap.String("trace_id", "123"), field.TraceID("123"))
	assert.Equal(t, zap.ByteString("user_agent", []byte("123")), field.UserAgent([]byte("123")))
	assert.Equal(t, zap.Strings("tags", []string{"123", "234"}), field.Tags([]string{"123", "234"}))
	assert.Equal(t, zap.Error(nil), field.Error(nil))
	assert.Equal(t, zap.Any("panic", nil), field.Panic(nil))
	assert.Equal(t, zap.Any("stack", nil), field.Stack(nil))
}

func Test_NewNop(t *testing.T) {
	assert.NotNil(t, NewNop())
}
