package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
	"sync"
)

// Level represents log levels
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Fields represents structured log fields
type Fields map[string]interface{}

// Entry represents a log entry
type Entry struct {
	Level     Level                 `json:"level"`
	Message   string                `json:"message"`
	Timestamp time.Time             `json:"timestamp"`
	Fields    Fields                `json:"fields,omitempty"`
	Caller    string                `json:"caller,omitempty"`
	RequestID string                `json:"request_id,omitempty"`
	UserID    string                `json:"user_id,omitempty"`
	Duration  time.Duration         `json:"duration,omitempty"`
	Error     error                 `json:"error,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Logger provides structured logging capabilities
type Logger struct {
	level    Level
	output   io.Writer
	fields   Fields
	mu       sync.RWMutex
	handlers []LogHandler
}

// LogHandler processes log entries
type LogHandler interface {
	Handle(entry Entry) error
}

// JSONHandler outputs logs in JSON format
type JSONHandler struct {
	output io.Writer
}

// NewJSONHandler creates a new JSON handler
func NewJSONHandler(output io.Writer) *JSONHandler {
	return &JSONHandler{output: output}
}

// Handle implements LogHandler interface
func (h *JSONHandler) Handle(entry Entry) error {
	// In a real implementation, you'd use a JSON encoder
	// For now, we'll use a simple format
	logLine := fmt.Sprintf(
		`{"level":"%s","message":"%s","timestamp":"%s","caller":"%s"}`,
		entry.Level.String(),
		entry.Message,
		entry.Timestamp.Format(time.RFC3339),
		entry.Caller,
	)
	
	if len(entry.Fields) > 0 {
		logLine += fmt.Sprintf(`,"fields":%v`, entry.Fields)
	}
	
	logLine += "\n"
	
	_, err := h.output.Write([]byte(logLine))
	return err
}

// NewLogger creates a new logger
func NewLogger(level Level, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	logger := &Logger{
		level:  level,
		output: output,
		fields: make(Fields),
	}

	// Add default JSON handler
	logger.AddHandler(NewJSONHandler(output))

	return logger
}

// AddHandler adds a log handler
func (l *Logger) AddHandler(handler LogHandler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers = append(l.handlers, handler)
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:    l.level,
		output:   l.output,
		fields:   make(Fields),
		handlers: l.handlers,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value

	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields Fields) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:    l.level,
		output:   l.output,
		fields:   make(Fields),
		handlers: l.handlers,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithContext adds context information to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l

	// Extract request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		logger = logger.WithField("request_id", requestID)
	}

	// Extract user ID from context
	if userID, ok := ctx.Value("user_id").(string); ok {
		logger = logger.WithField("user_id", userID)
	}

	return logger
}

// log creates and processes a log entry
func (l *Logger) log(level Level, message string, fields Fields) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	// Create entry
	entry := Entry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Fields:    fields,
		Caller:    caller,
	}

	// Add logger fields
	l.mu.RLock()
	for k, v := range l.fields {
		if entry.Fields == nil {
			entry.Fields = make(Fields)
		}
		entry.Fields[k] = v
	}
	l.mu.RUnlock()

	// Process through handlers
	l.mu.RLock()
	handlers := l.handlers
	l.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(entry); err != nil {
			// Fallback to standard log if handler fails
			log.Printf("Logger handler error: %v", err)
		}
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(DEBUG, message, f)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, err error, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = make(Fields)
	}

	if err != nil {
		f["error"] = err.Error()
	}

	l.log(ERROR, message, f)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields ...Fields) {
	var f Fields
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, message, f)
	os.Exit(1)
}

// RequestLogger logs HTTP request information
type RequestLogger struct {
	logger *Logger
}

// NewRequestLogger creates a new request logger
func NewRequestLogger(logger *Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// LogRequest logs an HTTP request
func (rl *RequestLogger) LogRequest(ctx context.Context, method, path, remoteAddr string, statusCode int, duration time.Duration) {
	fields := Fields{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"status_code": statusCode,
		"duration":    duration.String(),
	}

	level := INFO
	if statusCode >= 400 {
		level = WARN
	}
	if statusCode >= 500 {
		level = ERROR
	}

	rl.logger.WithContext(ctx).log(level, "HTTP Request", fields)
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(level Level) {
	globalLogger = NewLogger(level, os.Stdout)
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if globalLogger == nil {
		InitGlobalLogger(INFO)
	}
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(message string, fields ...Fields) {
	GetLogger().Debug(message, fields...)
}

// Info logs an info message using the global logger
func Info(message string, fields ...Fields) {
	GetLogger().Info(message, fields...)
}

// Warn logs a warning message using the global logger
func Warn(message string, fields ...Fields) {
	GetLogger().Warn(message, fields...)
}

// Error logs an error message using the global logger
func Error(message string, err error, fields ...Fields) {
	GetLogger().Error(message, err, fields...)
}

// Fatal logs a fatal message using the global logger
func Fatal(message string, fields ...Fields) {
	GetLogger().Fatal(message, fields...)
} 