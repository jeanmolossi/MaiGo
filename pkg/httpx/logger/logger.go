package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

var (
	_ (Logger) = (*noopLogger)(nil)
	_ (Logger) = (*jsonLogger)(nil)
	_ (Logger) = (*consoleLogger)(nil)
)

type noopLogger struct{}

func NewNoop() *noopLogger {
	return &noopLogger{}
}

// Error implements Logger.
func (n *noopLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	// noop
}

// Info implements Logger.
func (n *noopLogger) Info(ctx context.Context, msg string, args ...any) {
	// noop
}

type jsonLogger struct {
	log *slog.Logger
}

func NewJSON() *jsonLogger {
	return &jsonLogger{
		log: slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	}
}

// Error implements Logger.
func (j *jsonLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	j.log.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// Info implements Logger.
func (j *jsonLogger) Info(ctx context.Context, msg string, args ...any) {
	j.log.InfoContext(ctx, msg, args...)
}

type consoleLogger struct {
	log *slog.Logger
}

// Error implements Logger.
func (c *consoleLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	c.log.ErrorContext(ctx, msg, append([]any{"error", err}, args...)...)
}

// Info implements Logger.
func (c *consoleLogger) Info(ctx context.Context, msg string, args ...any) {
	c.log.InfoContext(ctx, msg, args...)
}

func NewConsole() *consoleLogger {
	return &consoleLogger{
		log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	}
}
