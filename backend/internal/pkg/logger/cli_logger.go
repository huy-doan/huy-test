package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type CLILoggerConfig struct {
	LogLevel         string
	LogDirectory     string
	EnableConsoleLog bool
}

// with files named cli-*.log
func InitCLILogger(config *CLILoggerConfig) Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	err = os.MkdirAll(config.LogDirectory, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	if config.EnableConsoleLog {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(os.Stderr)
	}

	return &loggerImpl{
		logger:       logger,
		traceID:      GenerateTraceID(),
		extraFields:  make(map[string]any),
		logDirectory: config.LogDirectory,
		filePrefix:   "cli",
	}
}

func (l *loggerImpl) getCliLogFile(level logrus.Level) *os.File {
	// Lock for file operations
	l.fileMutex.Lock()
	defer l.fileMutex.Unlock()

	var logPath string
	switch level {
	case logrus.DebugLevel:
		logPath = filepath.Join(l.logDirectory, "cli-debug.log")
	case logrus.InfoLevel:
		logPath = filepath.Join(l.logDirectory, "cli-info.log")
	case logrus.WarnLevel:
		logPath = filepath.Join(l.logDirectory, "cli-warning.log")
	case logrus.ErrorLevel:
		logPath = filepath.Join(l.logDirectory, "cli-error.log")
	default:
		logPath = filepath.Join(l.logDirectory, "cli-debug.log")
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return os.Stdout
	}

	return file
}

func (l *loggerImpl) logToCliFile(level logrus.Level, entry *logrus.Entry) {
	if !l.logger.IsLevelEnabled(level) {
		return
	}
	file := l.getCliLogFile(level)
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
