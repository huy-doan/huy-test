package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLLogger implements GORM's logger.Interface for database query logging
type SQLLogger struct {
	logger      *logrus.Logger
	level       logger.LogLevel
	traceLogger Logger
	config      *Config
}

// NewSQLLogger creates a new SQL logger that integrates with the application logger
func NewSQLLogger(config *Config, traceLogger Logger) logger.Interface {
	sqlLogger := logrus.New()

	// Set formatter to JSON for structured logging
	sqlLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Ensure log directory exists
	if err := os.MkdirAll(config.LogDirectory, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	// Configure DB log file
	dbLogPath := filepath.Join(config.LogDirectory, "db-backend.log")
	file, err := os.OpenFile(dbLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open db log file: %v", err))
	}

	// Set output destination
	sqlLogger.SetOutput(file)

	// Set SQL logging level based on app config
	sqlLogLevel := logger.Silent
	if config.EnableSQLLog {
		switch config.LogLevel {
		case "debug":
			sqlLogLevel = logger.Info
		case "info":
			sqlLogLevel = logger.Info
		case "warn":
			sqlLogLevel = logger.Warn
		case "error":
			sqlLogLevel = logger.Error
		}
	}

	return &SQLLogger{
		logger:      sqlLogger,
		level:       sqlLogLevel,
		traceLogger: traceLogger,
		config:      config,
	}
}

// LogMode sets the log level for SQL logger
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info logs SQL-related information messages
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.logger.WithFields(logrus.Fields{
			"trace_id": l.traceLogger.GetTraceID(),
			"type":     "sql_info",
		}).Info(fmt.Sprintf(msg, data...))
	}
}

// Warn logs SQL-related warning messages
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.logger.WithFields(logrus.Fields{
			"trace_id": l.traceLogger.GetTraceID(),
			"type":     "sql_warn",
		}).Warn(fmt.Sprintf(msg, data...))
	}
}

// Error logs SQL-related error messages
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.logger.WithFields(logrus.Fields{
			"trace_id": l.traceLogger.GetTraceID(),
			"type":     "sql_error",
		}).Error(fmt.Sprintf(msg, data...))
	}
}

// Trace logs SQL queries with execution details
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	// Calculate query execution time
	elapsed := time.Since(begin)

	// Get SQL query and affected rows
	sql, rows := fc()

	// Prepare log fields
	fields := logrus.Fields{
		"trace_id":   l.traceLogger.GetTraceID(),
		"elapsed_ms": elapsed.Milliseconds(),
		"rows":       rows,
		"sql":        sql,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	// If query resulted in an error
	if err != nil && err != gorm.ErrRecordNotFound {
		fields["error"] = err.Error()
		l.logger.WithFields(fields).Error("SQL Query Error")
		return
	}

	// Log slow queries as warnings
	if elapsed > time.Millisecond*500 && l.level >= logger.Warn {
		fields["is_slow"] = true
		l.logger.WithFields(fields).Warn("Slow SQL Query")
		return
	}

	// Log normal queries as info
	if l.level >= logger.Info {
		l.logger.WithFields(fields).Info("SQL Query")
	}
}
