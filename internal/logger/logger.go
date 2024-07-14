package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/imperiuse/go-app-skeleton/internal/metrics"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/imperiuse/go-app-skeleton/internal/config"
)

type (
	// Config - logger config.
	Config struct {
		Level    string
		Encoding string
		Color    bool
		Outputs  []string
		Tags     []string
	}

	// Logger logger wrapper under zap.Logger.
	Logger = zap.Logger
)

// NewNop - new Nop Logger.
func NewNop() *Logger {
	return zap.NewNop()
}

// New - create new logger.
func New(loggerCfg Config, e, serviceName, version string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	var lvl zapcore.Level

	err := lvl.UnmarshalText([]byte(loggerCfg.Level))
	if err != nil {
		return nil, fmt.Errorf("lvl.UnmarshalText: %w", err)
	}
	if e == config.Development {
		_ = lvl.UnmarshalText([]byte("debug"))
		loggerCfg.Color = true
		loggerCfg.Encoding = "console"
		loggerCfg.Outputs = []string{"stdout"}
	}

	cfg.Level.SetLevel(lvl)
	cfg.DisableStacktrace = true
	cfg.Development = e == config.Development
	cfg.Sampling.Initial = 50
	cfg.Sampling.Thereafter = 50
	cfg.Encoding = loggerCfg.Encoding
	cfg.OutputPaths = loggerCfg.Outputs
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if loggerCfg.Color {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("cfg.Build: %w", err)
	}

	logger = logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level >= zap.WarnLevel {
			metrics.LogsInc(entry.Level.String(), entry.Message)
		}

		return nil
	}))

	return logger.With(
			zap.String("env", e),
			zap.String("version", version),
			zap.String("services", serviceName),
		),
		nil
}

// LogIfError - log only if err!=nil.
func LogIfError(l *Logger, msg string, err error, f ...zapcore.Field) {
	LogCustomIfError(l.Error, msg, err, f...)
}

// LogCustomIfError - log with custom lvl only if err!=nil.
func LogCustomIfError(logFunc func(string, ...zapcore.Field), msg string, err error, f ...zapcore.Field) {
	if err != nil {
		logFunc(msg, append(f, zap.Error(err))...)
	}
}

type LoggerForES struct {
	*Logger
}

func NewLoggerForEs(log *Logger) *LoggerForES {
	return &LoggerForES{log}
}
func (l *LoggerForES) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	l.Sugar().Infof("[LoggerForES] LogRoundTrip %s %s %s [status:%d request:%s] err: %v",
		start.Format(time.RFC3339),
		req.Method,
		req.URL.String(),
		resStatusCode(res),
		dur.Truncate(time.Millisecond),
		err,
	)
	return nil
}
func (l *LoggerForES) RequestBodyEnabled() bool {
	return true
}
func (l *LoggerForES) ResponseBodyEnabled() bool {
	return true
}
func resStatusCode(res *http.Response) int {
	if res == nil {
		return -1
	}
	return res.StatusCode
}
