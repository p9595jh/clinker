package dbconn

import (
	"clinker-backend/common/config"
	"clinker-backend/common/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

type customLogger struct{}

func (*customLogger) name() string {
	return "DATABASE"
}

func (l *customLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	return &customLogger{}
}

func (l *customLogger) Info(ctx context.Context, s string, i ...interface{}) {}

func (l *customLogger) Warn(ctx context.Context, s string, i ...interface{}) {}

func (l *customLogger) Error(ctx context.Context, s string, i ...interface{}) {
	logger.Error(l.name()).D("data", i).W(s)
}

func (l *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if config.V.GetBool("db.log") {
		sql, rowsAffected := fc()
		fmt.Printf("%s [rowsAffected: %d]\n", sql, rowsAffected)
	} else {
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				sql, _ := fc()
				logger.Error(l.name()).E(err).D("sql", sql).W()
			}
		}
	}
}
