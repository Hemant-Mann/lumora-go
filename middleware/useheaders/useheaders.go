package useheaders

import (
	"bytes"
	"encoding/json"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
	"github.com/hemant-mann/lumora-go/core"
)

// SchemaWithParse is an interface for schemas that have a Parse method
// This allows us to work with different schema types
type SchemaWithParse interface {
	Parse(data any, destPtr any, options ...z.ExecOption) z.ZogIssueList
}

// UseHeaders creates a middleware that parses and validates request headers
// Similar to Lumora JS useHeaders hook
// schema: zog schema for validation (e.g., z.Struct(z.Shape{...}))
// dest: pointer to struct that will hold the parsed headers
func UseHeaders(schema SchemaWithParse, dest any) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) (*core.Response, error) {
			// Extract headers from request
			req := ctx.Request()
			headers := make(map[string]string)

			// Convert headers to map for zog parsing
			// Normalize header names to lowercase for consistent matching
			for name, values := range req.Header {
				// HTTP headers can have multiple values, we'll take the first one
				if len(values) > 0 {
					// Normalize to lowercase for zog schema matching
					normalizedName := strings.ToLower(name)
					headers[normalizedName] = values[0]
				}
			}

			jsonBody, err := json.Marshal(headers)
			if err != nil {
				return nil, core.NewError(400, "Failed to marshal headers to JSON")
			}

			// Validate headers using zog schema
			issues := schema.Parse(zjson.Decode(bytes.NewReader(jsonBody)), dest)
			if len(issues) > 0 {
				// Return validation error response
				resp := core.NewResponse().
					WithStatus(400).
					WithBody(map[string]string{"error": formatValidationErrors(issues)})
				return resp, nil
			}

			// Store parsed headers in context for access
			ctx.Set("_headers", dest)

			return next(ctx)
		}
	}
}

// UseHeadersWithKey creates a middleware that parses headers and stores them with a custom key
func UseHeadersWithKey(schema SchemaWithParse, dest any, key string) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) (*core.Response, error) {
			// Extract headers from request
			req := ctx.Request()
			headers := make(map[string]any)

			// Convert headers to map for zog parsing
			// Normalize header names to lowercase for consistent matching
			for name, values := range req.Header {
				// HTTP headers can have multiple values, we'll take the first one
				if len(values) > 0 {
					// Normalize to lowercase for zog schema matching
					normalizedName := strings.ToLower(name)
					headers[normalizedName] = values[0]
				}
			}
			jsonBody, err := json.Marshal(headers)
			if err != nil {
				return nil, core.NewError(400, "Failed to marshal headers to JSON")
			}

			// Validate headers using zog schema
			issues := schema.Parse(zjson.Decode(bytes.NewReader(jsonBody)), dest)
			if len(issues) > 0 {
				// Return validation error response
				resp := core.NewResponse().
					WithStatus(400).
					WithBody(map[string]string{"error": formatValidationErrors(issues)})
				return resp, nil
			}

			// Store parsed headers in context with custom key
			ctx.Set(key, dest)

			return next(ctx)
		}
	}
}

// GetHeaders retrieves the parsed headers from context
func GetHeaders(ctx core.Context) any {
	if headers, ok := ctx.Get("_headers"); ok {
		return headers
	}
	return nil
}

// formatValidationErrors formats zog validation errors into a readable string
func formatValidationErrors(issues z.ZogIssueList) string {
	if len(issues) == 0 {
		return "Validation failed"
	}

	// Format errors as a simple message
	// Could be enhanced to return structured error response
	msg := "Validation errors: "
	for i := 0; i < len(issues); i++ {
		if i > 0 {
			msg += "; "
		}
		if issues[i] != nil {
			msg += issues[i].Error()
		}
	}
	return msg
}
