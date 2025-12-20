package errorhandler

import (
	"github.com/hemant-mann/lumora-go/core"
)

// Options represents error handler configuration options
type Options struct {
	// Handler is a custom error handler function that returns a Response
	Handler func(ctx core.Context, err error) (*core.Response, error)
	// LogErrors determines if errors should be logged
	LogErrors bool
}

// DefaultOptions returns default error handler options
func DefaultOptions() *Options {
	return &Options{
		Handler:   defaultErrorHandler,
		LogErrors: true,
	}
}

// New creates a new error handler middleware
func New(options *Options) core.Middleware {
	if options == nil {
		options = DefaultOptions()
	}
	
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) (*core.Response, error) {
			resp, err := next(ctx)
			
			if err != nil {
				return options.Handler(ctx, err)
			}
			
			return resp, nil
		}
	}
}

// defaultErrorHandler handles errors by checking if they're HTTP errors
func defaultErrorHandler(ctx core.Context, err error) (*core.Response, error) {
	// Check if it's an HTTP error
	if httpErr := core.GetHTTPError(err); httpErr != nil {
		// Return error response
		resp := core.NewResponse().
			WithStatus(httpErr.Code).
			WithBody(map[string]string{"error": httpErr.Message})
		return resp, nil
	}
	
	// Default to 500 Internal Server Error
	resp := core.NewResponse().
		WithStatus(500).
		WithBody(map[string]string{"error": "Internal Server Error"})
	return resp, nil
}

// Simple creates a simple error handler middleware with default options
func Simple() core.Middleware {
	return New(DefaultOptions())
}

