package logger

// Logger is the custom structured logger interface
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)

	// Error logging extensions
	LogError(err error, contextData map[string]any)
	LogErrorWithContext(err error, msg string, contextData map[string]any)
	ErrorWithContext(format string, args ...any)

	// Context management
	WithTraceID(traceID string) Logger
	GetTraceID() string
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}
