package logging

import (
	"log"
	"time"

	"github.com/hemant-mann/lumora-go/core"
)

// Logger interface for custom loggers
type Logger interface {
	Log(level, message string, fields map[string]any)
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct{}

func (l *DefaultLogger) Log(level, message string, fields map[string]any) {
	log.Printf("[%s] %s %v", level, message, fields)
}

// Options represents logging configuration options
type Options struct {
	Logger Logger
	Format string // "json" or "text"
}

// DefaultOptions returns default logging options
func DefaultOptions() *Options {
	return &Options{
		Logger: &DefaultLogger{},
		Format: "text",
	}
}

// New creates a new logging middleware
func New(options *Options) core.Middleware {
	if options == nil {
		options = DefaultOptions()
	}

	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) (*core.Response, error) {
			start := time.Now()
			req := ctx.Request()

			// Log request
			logRequest(options, req.Method, req.URL.Path, req.RemoteAddr)

			// Execute next handler
			resp, err := next(ctx)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			statusCode := 200
			if resp != nil {
				statusCode = resp.StatusCode
			}

			logResponse(options, req.Method, req.URL.Path, statusCode, duration, err)

			return resp, err
		}
	}
}

func logRequest(options *Options, method, path, remoteAddr string) {
	fields := map[string]any{
		"method": method,
		"path":   path,
		"remote": remoteAddr,
	}
	options.Logger.Log("INFO", "Request started", fields)
}

func logResponse(options *Options, method, path string, statusCode int, duration time.Duration, err error) {
	fields := map[string]any{
		"method":      method,
		"path":        path,
		"status":      statusCode,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		fields["error"] = err.Error()
		options.Logger.Log("ERROR", "Request failed", fields)
	} else {
		options.Logger.Log("INFO", "Request completed", fields)
	}
}

// Simple creates a simple logging middleware with default options
func Simple() core.Middleware {
	return New(DefaultOptions())
}
