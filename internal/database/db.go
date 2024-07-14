package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	gormLogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB = gorm.DB

// DB - wrapper over *gorm.DB.
type DB struct {
	*gorm.DB

	log gormLogger.Interface

	shutdowner fx.Shutdowner
	errCounter atomic.Int32
}

const (
	maxOpenCons     = 20
	maxIdleCons     = 10
	maxConnLifeTime = time.Hour
	maxErrCount     = 5
)

// New - create new DB.
func New(dsn string, isDisableAllConstraints bool, gLogger gormLogger.Interface, shutdowner fx.Shutdowner) (*DB, error) {
	gormDB, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			SkipDefaultTransaction:                   true,
			DisableForeignKeyConstraintWhenMigrating: isDisableAllConstraints,
			Logger:                                   gLogger,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("gorm.Open error: %w", err)
	}

	db := &DB{DB: gormDB, log: gLogger, errCounter: atomic.Int32{}, shutdowner: shutdowner}

	sqlDB, err := db.GetSQLDB()
	if err != nil {
		return nil, fmt.Errorf("db.GetSQLDB error: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleCons)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenCons)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)

	return db, nil
}

// NewWithoutFX - simple way for create DB instance, use for small scripts and test purposes.
func NewWithoutFX(dsn string, isDisableAllConstraints bool, gLogger gormLogger.Interface) (*DB, error) {
	return New(dsn, isDisableAllConstraints, gLogger, nil)
}

// SetCustomGormObj - set custom gorm obj.
func (d *DB) SetCustomGormObj(g *gorm.DB) *DB {
	d.DB = g

	return d
}

// GetSQLDB - get generic database interface *sql.DB.
func (d *DB) GetSQLDB() (*sql.DB, error) {
	sqlDB, err := d.DB.DB()

	if err != nil {
		return nil, fmt.Errorf(
			"could not  generic database interface *sql.DB from the current *gorm.DB error: %w", err)
	}

	return sqlDB, nil
}

// IncreaseErrCnt -  increase err count, if count more than maxErrCount.
func (d *DB) IncreaseErrCnt() {
	d.errCounter.Add(1)

	if d.errCounter.Load() > maxErrCount {
		const desc = "db.IncreaseErrCnt max error count achieved, shutdown app => cancel global context."
		d.log.Error(context.Background(), desc)

		if d.shutdowner == nil {
			panic(desc)
		}

		if err := d.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
			panic(fmt.Errorf("%w due situation: %s", err, desc))
		}
	}
}

// FlushErrCnt -  set up zero err.
func (d *DB) FlushErrCnt() {
	d.errCounter.Store(0)
}

// Ping -  just ping.
func (d *DB) Ping() error {
	sqlDB, err := d.GetSQLDB()
	if err != nil {
		return fmt.Errorf("ping db error: %w", err)
	}

	return sqlDB.Ping()
}

// Stats - Returns database statistics.
func (d *DB) Stats() (*sql.DBStats, error) {
	sqlDB, err := d.GetSQLDB()
	if err != nil {
		return nil, fmt.Errorf("stats db error: %w", err)
	}

	stats := sqlDB.Stats()
	return &stats, nil
}

// Close - close conn.
func (d *DB) Close() error {
	sqlDB, err := d.GetSQLDB()
	if err != nil {
		return fmt.Errorf("close db error: %w", err)
	}

	return sqlDB.Close()
}
