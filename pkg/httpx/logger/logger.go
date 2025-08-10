package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger abstracts the logging backend used by the round tripper. Implementations
// must be safe for concurrent use.
type Logger interface {
	// Info writes an informational log message with optional attributes.
	Info(ctx context.Context, msg string, args ...any)
	// Error writes an error log message with optional attributes.
	Error(ctx context.Context, err error, msg string, args ...any)
}

var (
	_ (Logger) = (*noopLogger)(nil)
	_ (Logger) = (*jsonLogger)(nil)
	_ (Logger) = (*consoleLogger)(nil)
)

type noopLogger struct{}

// NewNoop returns a Logger that discards all logs.
func NewNoop() *noopLogger {
	return &noopLogger{}
}

// Error implements Logger but performs no operation.
func (n *noopLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	// noop
}

// Info implements Logger but performs no operation.
func (n *noopLogger) Info(ctx context.Context, msg string, args ...any) {
	// noop
}

type jsonLogger struct {
	log *slog.Logger
}

// NewJSON creates a Logger that emits structured JSON logs using slog.
func NewJSON() *jsonLogger {
	return &jsonLogger{
		log: slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	}
}

// Error implements Logger by logging the provided error in JSON format.
func (j *jsonLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	j.log.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// Info implements Logger by logging the message in JSON format.
func (j *jsonLogger) Info(ctx context.Context, msg string, args ...any) {
	j.log.InfoContext(ctx, msg, args...)
}

type consoleLogger struct {
	log *slog.Logger
}

// Error implements Logger by logging the provided error in text format.
func (c *consoleLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	c.log.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// Info implements Logger by logging the message in text format.
func (c *consoleLogger) Info(ctx context.Context, msg string, args ...any) {
	c.log.InfoContext(ctx, msg, args...)
}

// NewConsole creates a Logger that writes human-readable text logs to stdout.
func NewConsole() *consoleLogger {
	return &consoleLogger{
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	}
}
