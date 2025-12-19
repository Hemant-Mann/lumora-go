package cors

import (
	"fmt"

	"github.com/hemant-mann/lumora-go/core"
)

// Options represents CORS configuration options
type Options struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
	MaxAge         int
	AllowCredentials bool
}

// DefaultOptions returns default CORS options
func DefaultOptions() *Options {
	return &Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{},
		MaxAge:         86400,
		AllowCredentials: false,
	}
}

// New creates a new CORS middleware
func New(options *Options) core.Middleware {
	if options == nil {
		options = DefaultOptions()
	}
	
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) error {
			origin := ctx.Header("Origin")
			
			// Handle preflight request
			if ctx.Request().Method == "OPTIONS" {
				setCORSHeaders(ctx, options, origin)
				ctx.Status(204)
				return nil
			}
			
			// Set CORS headers for actual request
			setCORSHeaders(ctx, options, origin)
			
			return next(ctx)
		}
	}
}

func setCORSHeaders(ctx core.Context, options *Options, origin string) {
	// Set Access-Control-Allow-Origin
	if len(options.AllowedOrigins) > 0 {
		if options.AllowedOrigins[0] == "*" {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
		} else {
			// Check if origin is in allowed list
			for _, allowedOrigin := range options.AllowedOrigins {
				if allowedOrigin == origin {
					ctx.SetHeader("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
	}
	
	// Set Access-Control-Allow-Methods
	if len(options.AllowedMethods) > 0 {
		methods := ""
		for i, method := range options.AllowedMethods {
			if i > 0 {
				methods += ", "
			}
			methods += method
		}
		ctx.SetHeader("Access-Control-Allow-Methods", methods)
	}
	
	// Set Access-Control-Allow-Headers
	if len(options.AllowedHeaders) > 0 {
		headers := ""
		for i, header := range options.AllowedHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += header
		}
		ctx.SetHeader("Access-Control-Allow-Headers", headers)
	}
	
	// Set Access-Control-Expose-Headers
	if len(options.ExposedHeaders) > 0 {
		headers := ""
		for i, header := range options.ExposedHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += header
		}
		ctx.SetHeader("Access-Control-Expose-Headers", headers)
	}
	
	// Set Access-Control-Max-Age
	if options.MaxAge > 0 {
		ctx.SetHeader("Access-Control-Max-Age", fmt.Sprintf("%d", options.MaxAge))
	}
	
	// Set Access-Control-Allow-Credentials
	if options.AllowCredentials {
		ctx.SetHeader("Access-Control-Allow-Credentials", "true")
	}
}

