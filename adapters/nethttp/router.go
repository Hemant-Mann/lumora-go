package nethttp

import (
	"strings"

	"github.com/hemant-mann/lumora-go/core"
)

type route struct {
	method  string
	pattern string
	handler func(core.Context) error // Wrapper that returns error for router compatibility
}

type Router struct {
	routes []*route
}

func NewRouter() *Router {
	return &Router{
		routes: []*route{},
	}
}

func (r *Router) Handle(method, pattern string, handler func(core.Context) error) {
	r.routes = append(r.routes, &route{
		method:  method,
		pattern: pattern,
		handler: handler,
	})
}

func (r *Router) Match(method, path string) (func(core.Context) error, map[string]string) {
	for _, route := range r.routes {
		if route.method != method {
			continue
		}
		
		params := matchPattern(route.pattern, path)
		if params != nil {
			return route.handler, params
		}
	}
	
	return nil, nil
}

// matchPattern matches a pattern like "/users/:id" against a path like "/users/123"
// Returns nil if no match, otherwise returns a map of parameter names to values
func matchPattern(pattern, path string) map[string]string {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	
	if len(patternParts) != len(pathParts) {
		return nil
	}
	
	params := make(map[string]string)
	
	for i, patternPart := range patternParts {
		pathPart := pathParts[i]
		
		if strings.HasPrefix(patternPart, ":") {
			// It's a parameter
			paramName := strings.TrimPrefix(patternPart, ":")
			params[paramName] = pathPart
		} else if patternPart != pathPart {
			// Literal mismatch
			return nil
		}
	}
	
	return params
}

