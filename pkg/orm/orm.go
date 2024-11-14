package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

type DB struct {
	*gorm.DB
}

type OrmLog struct {
	logLevel logger.LogLevel
}

func (l *OrmLog) LogMode(logLevel logger.LogLevel) logger.Interface {
	l.logLevel = logLevel
	return l
}

func (l *OrmLog) Info(ctx context.Context, format string, args ...any) {
	if l.logLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(format, args...)
}

func (l *OrmLog) Warn(ctx context.Context, format string, args ...any) {
	if l.logLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Infof(format, args...)
}

func (l *OrmLog) Error(ctx context.Context, format string, args ...any) {
	if l.logLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(format, args...)
}

func (l *OrmLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}

func newDB(c *Config) (*DB, error) {
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 10
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 100
	}
	if c.MaxLifetime == 0 {
		c.MaxLifetime = 3600
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: &OrmLog{},
	})
	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxOpenConns(c.MaxOpenConns)
	sdb.SetMaxIdleConns(c.MaxIdleConns)
	sdb.SetConnMaxLifetime(time.Duration(c.MaxLifetime))
	err = db.Use(NewCustomePlugin())
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func MustNewMysql(c *Config) *DB {
	db, err := newDB(c)
	if err != nil {
		panic(err)
	}
	return db
}
