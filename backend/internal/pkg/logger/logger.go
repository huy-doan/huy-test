package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"maps"

	"github.com/google/uuid"
	appConfig "github.com/huydq/test/internal/pkg/config"
	"github.com/sirupsen/logrus"
)

// LogLevel represents logging levels
type LogLevel string

const (
	// DEBUG level for detailed information in development environment
	DEBUG LogLevel = "debug"
	// INFO level for general operational information
	INFO LogLevel = "info"
	// WARN level for warnings that don't cause errors but should be noted
	WARN LogLevel = "warn"
	// ERROR level for system errors
	ERROR LogLevel = "error"
)

// TraceIDKey is the context key for trace ID
const TraceIDKey = "trace_id"

// loggerImpl is the implementation of Logger interface
type loggerImpl struct {
	logger       *logrus.Logger
	traceID      string
	extraFields  map[string]any
	fileMutex    sync.Mutex
	logDirectory string
	filePrefix   string // Prefix for log files, empty for web, "cli" for CLI
}

// Config holds the configuration for logger
type Config struct {
	LogLevel         string
	LogDirectory     string
	EnableConsoleLog bool
	EnableSQLLog     bool
}

// Global singleton instance and initialization lock
var (
	instance      Logger
	instanceOnce  sync.Once
	instanceMutex sync.RWMutex
)

// InitLogger initializes the global logger instance with the given config
// This should be called once at application startup
func InitLogger(config *Config) {
	instanceOnce.Do(func() {
		logger := logrus.New()

		level, err := logrus.ParseLevel(config.LogLevel)
		if err != nil {
			level = logrus.InfoLevel
		}
		logger.SetLevel(level)

		err = os.MkdirAll(config.LogDirectory, 0755)
		if err != nil {
			// If we can't create the directory, log to stderr and keep going
			logger.WithField("error", err.Error()).Error("Failed to create log directory, logging to stderr")
		}

		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})

		if config.EnableConsoleLog {
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(io.Discard)
		}

		instance = &loggerImpl{
			logger:       logger,
			traceID:      GenerateTraceID(),
			extraFields:  make(map[string]any),
			logDirectory: config.LogDirectory,
		}
	})
}

// GetLogger returns the global logger instance
// If the logger hasn't been initialized, it returns a default logger
func GetLogger() Logger {
	instanceMutex.RLock()
	defer instanceMutex.RUnlock()

	if instance == nil {
		appConfig := appConfig.LoadConfig()

		defaultConfig := &Config{
			LogLevel:         appConfig.LogLevel,
			LogDirectory:     appConfig.LogDirectory,
			EnableConsoleLog: appConfig.EnableConsoleLog,
			EnableSQLLog:     appConfig.EnableSQLLog,
		}

		return NewLogger(defaultConfig)
	}

	return instance
}

// NewLogger creates a new logger instance
// This should be used only for specific cases where a separate logger instance is needed
func NewLogger(config *Config) Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	err = os.MkdirAll(config.LogDirectory, 0755)
	if err != nil {
		// If we can't create the directory, log to stderr and keep going
		logger.WithField("error", err.Error()).Error("Failed to create log directory, logging to stderr")
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	if config.EnableConsoleLog {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(io.Discard) // Discard default output as we'll use level-specific files
	}

	return &loggerImpl{
		logger:       logger,
		traceID:      GenerateTraceID(),
		extraFields:  make(map[string]any),
		logDirectory: config.LogDirectory,
	}
}

// GenerateTraceID generates a unique trace ID for request tracking
func GenerateTraceID() string {
	return uuid.New().String()
}

// WithTraceID creates a new logger instance with the specified trace ID
func (l *loggerImpl) WithTraceID(traceID string) Logger {
	newLogger := &loggerImpl{
		logger:       l.logger,
		traceID:      traceID,
		extraFields:  make(map[string]any),
		logDirectory: l.logDirectory,
		filePrefix:   l.filePrefix,
	}

	// Copy extra fields
	maps.Copy(newLogger.extraFields, l.extraFields)

	return newLogger
}

// WithField adds a field to the logger
func (l *loggerImpl) WithField(key string, value any) Logger {
	newLogger := &loggerImpl{
		logger:       l.logger,
		traceID:      l.traceID,
		extraFields:  make(map[string]any),
		logDirectory: l.logDirectory,
		filePrefix:   l.filePrefix,
	}

	// Copy existing extra fields
	maps.Copy(newLogger.extraFields, l.extraFields)

	// Add new field
	newLogger.extraFields[key] = value

	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *loggerImpl) WithFields(fields map[string]any) Logger {
	newLogger := &loggerImpl{
		logger:       l.logger,
		traceID:      l.traceID,
		extraFields:  make(map[string]any),
		logDirectory: l.logDirectory,
		filePrefix:   l.filePrefix,
	}

	// Copy existing extra fields
	maps.Copy(newLogger.extraFields, l.extraFields)

	// Add new fields
	maps.Copy(newLogger.extraFields, fields)

	return newLogger
}

// GetTraceID returns the current trace ID
func (l *loggerImpl) GetTraceID() string {
	return l.traceID
}

// makeFields adds common fields to all log entries
func (l *loggerImpl) makeFields(fields map[string]any) logrus.Fields {
	if fields == nil {
		fields = make(map[string]any)
	}

	// Add trace ID and timestamp
	fields[TraceIDKey] = l.traceID
	fields["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Add extra fields
	maps.Copy(fields, l.extraFields)

	return logrus.Fields(fields)
}

// getLogFile returns the appropriate log file for the given level
func (l *loggerImpl) getLogFile(level logrus.Level) *os.File {
	// Lock for file operations
	l.fileMutex.Lock()
	defer l.fileMutex.Unlock()

	// Use the stored log directory
	var logPath string
	switch level {
	case logrus.DebugLevel:
		logPath = filepath.Join(l.logDirectory, "debug.log")
	case logrus.InfoLevel:
		logPath = filepath.Join(l.logDirectory, "info.log")
	case logrus.WarnLevel:
		logPath = filepath.Join(l.logDirectory, "warning.log")
	case logrus.ErrorLevel:
		logPath = filepath.Join(l.logDirectory, "error.log")
	default:
		logPath = filepath.Join(l.logDirectory, "app.log")
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.logger.WithField("error", err.Error()).Error("Failed to open log file")
		return os.Stdout
	}

	return file
}

// logToFile logs a message to the appropriate file based on level
func (l *loggerImpl) logToFile(level logrus.Level, entry *logrus.Entry) {
	if !l.logger.IsLevelEnabled(level) {
		return
	}
	file := l.getLogFile(level)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing log file: %v\n", err)
		}
	}()

	fileLogger := logrus.New()
	fileLogger.SetOutput(file)
	fileLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	fileLogger.SetLevel(level)

	switch level {
	case logrus.DebugLevel:
		fileLogger.WithFields(entry.Data).Debug(entry.Message)
	case logrus.InfoLevel:
		fileLogger.WithFields(entry.Data).Info(entry.Message)
	case logrus.WarnLevel:
		fileLogger.WithFields(entry.Data).Warn(entry.Message)
	case logrus.ErrorLevel:
		fileLogger.WithFields(entry.Data).Error(entry.Message)
	}
}

// logToCliFile logs a message to the CLI file
func (l *loggerImpl) log(level logrus.Level, entry *logrus.Entry) {
	if l.filePrefix == "cli" {
		l.logToCliFile(level, entry)
	} else {
		l.logToFile(level, entry)
	}
}

// Debug logs a message at the DEBUG level
func (l *loggerImpl) Debug(msg string, fields map[string]any) {
	if l.logger.IsLevelEnabled(logrus.DebugLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.DebugLevel,
			Message: msg,
		}

		l.logger.WithFields(entry.Data).Debug(msg)

		l.log(logrus.DebugLevel, entry)
	}
}

// Info logs a message at the INFO level
func (l *loggerImpl) Info(msg string, fields map[string]any) {
	if l.logger.IsLevelEnabled(logrus.InfoLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.InfoLevel,
			Message: msg,
		}

		l.logger.WithFields(entry.Data).Info(msg)

		l.log(logrus.InfoLevel, entry)
	}
}

// Warn logs a message at the WARN level
func (l *loggerImpl) Warn(msg string, fields map[string]any) {
	if l.logger.IsLevelEnabled(logrus.WarnLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.WarnLevel,
			Message: msg,
		}

		l.logger.WithFields(entry.Data).Warn(msg)

		l.log(logrus.WarnLevel, entry)
	}
}

// Error logs a message at the ERROR level
func (l *loggerImpl) Error(msg string, fields map[string]any) {
	if l.logger.IsLevelEnabled(logrus.ErrorLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.ErrorLevel,
			Message: msg,
		}

		l.logger.WithFields(entry.Data).Error(msg)

		l.log(logrus.ErrorLevel, entry)
	}
}
