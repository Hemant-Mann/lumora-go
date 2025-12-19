package usejsonbody

import (
	"bytes"
	"io"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
	"github.com/hemant-mann/lumora-go/core"
)

// SchemaWithParse is an interface for schemas that have a Parse method
// This allows us to work with different schema types
type SchemaWithParse interface {
	Parse(data any, destPtr any, options ...z.ExecOption) z.ZogIssueList
}

// UseJsonBody creates a middleware that parses and validates JSON request body
// Similar to Lumora JS useJsonBody hook
// schema: zog schema for validation (e.g., z.Struct(z.Shape{...}))
// dest: pointer to struct that will hold the parsed data
func UseJsonBody(schema SchemaWithParse, dest interface{}) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) error {
			req := ctx.Request()

			// Read request body
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return core.NewError(400, "Failed to read request body")
			}
			defer req.Body.Close()

			if len(body) == 0 {
				return core.NewError(400, "Request body is empty")
			}

			// Decode and validate JSON using zjson
			issues := schema.Parse(zjson.Decode(bytes.NewReader(body)), dest)
			if len(issues) > 0 {
				// Return validation errors
				return core.NewError(400, formatValidationErrors(issues))
			}

			// Store parsed body in context for access
			ctx.Set("_jsonBody", dest)

			return next(ctx)
		}
	}
}

// UseJsonBodyWithKey creates a middleware that parses JSON and stores it with a custom key
func UseJsonBodyWithKey(schema SchemaWithParse, dest interface{}, key string) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) error {
			req := ctx.Request()

			// Read request body
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return core.NewError(400, "Failed to read request body")
			}
			defer req.Body.Close()

			// Decode and validate JSON using zjson
			issues := schema.Parse(zjson.Decode(bytes.NewReader(body)), dest)
			if len(issues) > 0 {
				// Return validation errors
				return core.NewError(400, formatValidationErrors(issues))
			}

			// Store parsed body in context with custom key
			ctx.Set(key, dest)

			return next(ctx)
		}
	}
}

// GetJsonBody retrieves the parsed JSON body from context
func GetJsonBody(ctx core.Context) any {
	if body, ok := ctx.Get("_jsonBody"); ok {
		return body
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
