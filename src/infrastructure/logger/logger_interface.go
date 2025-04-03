package logger

// Logger is the custom structured logger interface
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})

	// Error logging extensions
	LogError(err error, contextData map[string]interface{})
	LogErrorWithContext(err error, msg string, contextData map[string]interface{})
	ErrorWithContext(format string, args ...interface{})

	// Context management
	WithTraceID(traceID string) Logger
	GetTraceID() string
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}
