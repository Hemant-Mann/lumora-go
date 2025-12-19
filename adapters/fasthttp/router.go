package fasthttp

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Router wraps the fasthttp/router.Router
type Router struct {
	router *router.Router
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		router: router.New(),
	}
}

// Handle registers a handler for a method and path pattern
// Note: fasthttp/router uses {name} for parameters, not :name
// The handler parameter is a function that takes *fasthttp.RequestCtx and handles errors internally
func (r *Router) Handle(method, pattern string, handler fasthttp.RequestHandler) {
	// Convert :param to {param} format for fasthttp/router
	convertedPattern := convertPattern(pattern)

	// Register with the router based on method
	switch method {
	case "GET":
		r.router.GET(convertedPattern, handler)
	case "POST":
		r.router.POST(convertedPattern, handler)
	case "PUT":
		r.router.PUT(convertedPattern, handler)
	case "DELETE":
		r.router.DELETE(convertedPattern, handler)
	case "PATCH":
		r.router.PATCH(convertedPattern, handler)
	default:
		r.router.Handle(method, convertedPattern, handler)
	}
}

// Handler returns the fasthttp.RequestHandler
func (r *Router) Handler() fasthttp.RequestHandler {
	return r.router.Handler
}

// convertPattern converts :param format to {param} format
// Example: /users/:id -> /users/{id}
func convertPattern(pattern string) string {
	// Simple conversion: replace :param with {param}
	// This is a basic implementation - could be enhanced with regex for edge cases
	result := ""
	i := 0
	for i < len(pattern) {
		if i < len(pattern)-1 && pattern[i] == ':' && (i == 0 || pattern[i-1] == '/') {
			// Found :param, convert to {param}
			result += "{"
			i++ // Skip the ':'
			// Read the parameter name
			paramStart := i
			for i < len(pattern) && pattern[i] != '/' {
				i++
			}
			result += pattern[paramStart:i]
			result += "}"
		} else {
			result += string(pattern[i])
			i++
		}
	}
	return result
}
