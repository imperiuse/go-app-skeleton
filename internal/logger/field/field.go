package field

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// f package is used to log fields with strict types.
// Package prevents mapping conflicts in elastic.

func Bool(name string, b bool) zapcore.Field {
	return zap.Bool(name, b)
}

func Int(name string, i int) zapcore.Field {
	return zap.Int(name, i)
}

func Int64(name string, i int64) zapcore.Field {
	return zap.Int64(name, i)
}

func Uint64(name string, i uint64) zapcore.Field {
	return zap.Uint64(name, i)
}

func String(name string, s string) zapcore.Field {
	return zap.String(name, s)
}

func Any(name string, obj any) zapcore.Field {
	return zap.Any(name, obj)
}

func ID(idi any) zapcore.Field {
	switch id := idi.(type) {
	case int64:
		return zap.Int64("id", id)
	case int32:
		return zap.Int32("id", id)
	case int:
		return zap.Int("id", id)
	case string:
		return zap.String("id", id)
	default:
		return zap.String("id", fmt.Sprint(id))
	}
}

func Controller(name string) zapcore.Field {
	return zap.String("controller", name)
}

func Handler(name string) zapcore.Field {
	return zap.String("handler", name)
}

func Topic(topic string) zapcore.Field {
	return zap.String("topic", topic)
}

func Repo(name string) zapcore.Field {
	return zap.String("repo", name)
}

func Table(name string) zapcore.Field {
	return zap.String("table", name)
}

func Env(env string) zapcore.Field {
	return zap.String("env", env)
}

func Version(version string) zapcore.Field {
	return zap.String("version", version)
}

func Service(service string) zapcore.Field {
	return zap.String("service", service)
}

func Tags(tags []string) zapcore.Field {
	return zap.Strings("tags", tags)
}

func Port(port int) zapcore.Field {
	return zap.Int("port", port)
}

func Host(host string) zapcore.Field {
	return zap.String("host", host)
}

func Path(path string) zapcore.Field {
	return zap.String("path", path)
}

func Error(err error) zapcore.Field {
	return zap.Error(err)
}

func Panic(pan any) zapcore.Field {
	return zap.Any("panic", pan)
}

func Stack(stack any) zapcore.Field {
	return zap.Any("stack", stack)
}

func IP(ip string) zapcore.Field {
	return zap.String("ip", ip)
}

func IPForwardedFor(ip string) zapcore.Field {
	return zap.String("ip_forwarded_for", ip)
}

func TraceID(traceID string) zapcore.Field {
	return zap.String("trace_id", traceID)
}

func UserAgent(ua []byte) zapcore.Field {
	return zap.ByteString("user_agent", ua)
}
