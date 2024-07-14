package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLoggerConfig = gormLogger.Config

type GormLogger struct {
	log *zap.Logger
	cfg GormLoggerConfig
}

const (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

const mgn = 1e6

// NewGormLogger - create new custom logger for GORM.
func NewGormLogger(l *zap.Logger, config GormLoggerConfig) gormLogger.Interface {
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 100 * time.Millisecond
	}

	return &GormLogger{
		log: l,
		cfg: config,
	}
}

// LogMode - it's only needed for support Debug() method in gormLogger.
func (l *GormLogger) LogMode(lvl gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l // this is needed! it's important, deep copy is right here.
	newLogger.cfg.LogLevel = lvl

	return &newLogger
}

func (l *GormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Infof(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (l *GormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Warnf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (l *GormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Errorf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
}

func (l *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.cfg.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.cfg.LogLevel >= gormLogger.Error &&
		(!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.cfg.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.log.Sugar().Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/mgn, "-", sql)
		} else {
			l.log.Sugar().Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/mgn, rows, sql)
		}
	case elapsed > l.cfg.SlowThreshold && l.cfg.SlowThreshold != 0 && l.cfg.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.cfg.SlowThreshold)
		if rows == -1 {
			l.log.Sugar().Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/mgn, "-", sql)
		} else {
			l.log.Sugar().Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/mgn, rows, sql)
		}
	case l.cfg.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.log.Sugar().Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/mgn, "-", sql)
		} else {
			l.log.Sugar().Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/mgn, rows, sql)
		}
	}
}
