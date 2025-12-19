package core

import (
	"context"
	"net/http"
)

// Context represents the request context with framework-agnostic abstractions
type Context interface {
	// Request returns the underlying HTTP request
	Request() *http.Request
	
	// Response returns the response writer
	Response() http.ResponseWriter
	
	// Get retrieves a value from the context
	Get(key string) (interface{}, bool)
	
	// Set stores a value in the context
	Set(key string, value interface{})
	
	// Param returns a path parameter by name
	Param(name string) string
	
	// Query returns a query parameter by name
	Query(name string) string
	
	// Header returns a request header by name
	Header(name string) string
	
	// SetHeader sets a response header
	SetHeader(name, value string)
	
	// Status sets the HTTP status code
	Status(code int)
	
	// JSON sends a JSON response
	JSON(code int, data interface{}) error
	
	// String sends a text response
	String(code int, format string, values ...interface{}) error
	
	// BindJSON binds the request body to a struct
	BindJSON(dest interface{}) error
	
	// Context returns the underlying context.Context
	Context() context.Context
	
	// WithContext returns a new context with the given context.Context
	WithContext(ctx context.Context) Context
}

