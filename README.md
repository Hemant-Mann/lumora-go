# Lumora Go

A framework-agnostic, composable web framework for Go inspired by [Lumora](https://github.com/hemant-mann/lumora).

## Features

- **Framework Agnostic**: Easily switch between `net/http`, `Gin`, and `FastHTTP`
- **Composable API**: Clean functional composition pattern for middleware
- **Decoupled Design**: Modular packages for easy extension
- **Type-Safe**: Leverages Go's type system for safety

## Installation

```bash
# Clone the repository
git clone https://github.com/hemant-mann/lumora-go.git
cd lumora-go

# Install dependencies
go mod download

# Or if using as a dependency
go get github.com/hemant-mann/lumora-go
```

## Quick Start

### Using net/http

```go
package main

import (
	"github.com/hemant-mann/lumora-go/adapters/nethttp"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/middleware/cors"
	"github.com/hemant-mann/lumora-go/middleware/logging"
	"github.com/hemant-mann/lumora-go/middleware/errorhandler"
)

func main() {
	app := nethttp.New()
	
	// Add global middleware
	app.Use(
		cors.New(cors.DefaultOptions()),
		logging.Simple(),
		errorhandler.Simple(),
	)
	
	// Register services
	app.Services().Register("userService", NewUserService())
	
	// Define routes
	app.Get("/", func(ctx core.Context) (*core.Response, error) {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello, World!",
			})
		return resp, nil
	})
	
	app.Get("/users/:id", func(ctx core.Context) (*core.Response, error) {
		id := ctx.Param("id")
		
		// Access service
		userService := ctx.MustService("userService").(*UserService)
		
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"id": id,
			})
		return resp, nil
	})
	
	app.Start(":8080")
}
```

### Using Gin

```go
package main

import (
	"github.com/hemant-mann/lumora-go/adapters/gin"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/middleware/cors"
	"github.com/hemant-mann/lumora-go/middleware/logging"
	"github.com/hemant-mann/lumora-go/middleware/errorhandler"
)

func main() {
	app := gin.New()
	
	app.Use(
		cors.New(cors.DefaultOptions()),
		logging.Simple(),
		errorhandler.Simple(),
	)
	
	app.Get("/", func(ctx core.Context) (*core.Response, error) {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello from Gin!",
			})
		return resp, nil
	})
	
	app.Start(":8080")
}
```

### Using FastHTTP

```go
package main

import (
	"github.com/hemant-mann/lumora-go/adapters/fasthttp"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/middleware/cors"
	"github.com/hemant-mann/lumora-go/middleware/logging"
	"github.com/hemant-mann/lumora-go/middleware/errorhandler"
)

func main() {
	app := fasthttp.New()
	
	app.Use(
		cors.New(cors.DefaultOptions()),
		logging.Simple(),
		errorhandler.Simple(),
	)
	
	app.Get("/", func(ctx core.Context) (*core.Response, error) {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello from FastHTTP!",
			})
		return resp, nil
	})
	
	app.Start(":8080")
}
```

## Architecture

### Core Package

The `core` package defines framework-agnostic interfaces:
- `Context`: Request/response context abstraction
- `Handler`: Handler function type `func(Context) (*Response, error)`
- `Middleware`: Middleware function type
- `App`: Application interface

### Adapters

Framework-specific implementations:
- `adapters/nethttp`: Standard library adapter
- `adapters/gin`: Gin framework adapter
- `adapters/fasthttp`: FastHTTP adapter

### Middleware

Composable middleware packages:
- `middleware/cors`: CORS handling
- `middleware/logging`: Request/response logging
- `middleware/errorhandler`: Error handling

## Composing Middleware

Middleware can be composed using the `core.Compose` function:

```go
middleware := core.Compose(
	middleware1,
	middleware2,
	middleware3,
)

handler := middleware(myHandler)
```

Or applied directly:

```go
handler := core.Apply(
	myHandler,
	middleware1,
	middleware2,
	middleware3,
)
```

## Error Handling

Use the error handling utilities:

```go
// Create an HTTP error
err := core.NewError(404, "Not Found")

// Wrap an existing error
err := core.WrapError(500, "Internal Server Error", originalErr)

// In your handler - return error as second value
if someCondition {
	return nil, core.NewError(400, "Bad Request")
}
```

## Response System

The response system gives you full control over your response structure. Handlers return `(*Response, error)` instead of sending responses directly. The orchestrator handles sending the response automatically.

```go
// Create a response with full control
resp := core.NewResponse().
	WithStatus(200).
	WithHeader("X-Custom-Header", "value").
	WithCookie(core.Cookie{
		Name:     "session",
		Value:    "abc123",
		HttpOnly: true,
		MaxAge:   3600,
	}).
	WithBody(map[string]interface{}{
		"message": "Success",
		"data":    yourData,
	})
return resp, nil  // Orchestrator will send the response
```

### Automatic Content-Type Detection

- If `Body` is a `string`, it's sent as `text/plain`
- Otherwise, it's sent as `application/json`

```go
// Plain text response
resp := core.NewResponse().
	WithStatus(200).
	WithBody("This is plain text")

// JSON response (automatic)
resp := core.NewResponse().
	WithStatus(200).
	WithBody(map[string]string{"key": "value"})
```

### Helper Functions

```go
// Quick JSON response
core.JSON(ctx, 200, data)

// Quick text response
core.Text(ctx, 200, "Hello")

// Formatted string response
core.String(ctx, 200, "Hello, %s", name)
```

## Services (Dependency Injection)

Services provide dependency injection capabilities, similar to Lumora JS. You can register services at the app level or per-route.

### App-Level Services

Services registered at the app level are available to all routes:

```go
// Define a service
type UserService struct {
	users map[string]string
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]string),
	}
}

// Register app-level services
app.Services().Register("userService", NewUserService())
```

### Route-Specific Services (Like Lumora JS!)

You can compose routes with route-specific services using the `UseServices` middleware:

```go
import "github.com/hemant-mann/lumora-go/middleware/useservices"

// Route with route-specific services
app.Get("/users/:id",
	func(ctx core.Context) (*core.Response, error) {
		// Access route-specific service
		userService := ctx.MustService("userService").(*UserService)
		// ... use the service
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{"id": "123"})
		return resp, nil
	},
	useservices.UseServices(map[string]interface{}{
		"userService": NewUserService(),
	}),
)

// Multiple route-specific services
app.Get("/protected/:id",
	func(ctx core.Context) (*core.Response, error) {
		userService := ctx.MustService("userService").(*UserService)
		authService := ctx.MustService("authService").(*AuthService)
		// ... use services
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{"message": "Protected"})
		return resp, nil
	},
	useservices.UseServices(map[string]interface{}{
		"userService": NewUserService(),
		"authService": NewAuthService(), // Overrides app-level service
	}),
)

// Single service helper
app.Post("/users",
	func(ctx core.Context) (*core.Response, error) {
		// ... handler code
		resp := core.NewResponse().WithStatus(201)
		return resp, nil
	},
	useservices.UseService("userService", NewUserService()),
)
```

### Service Resolution

Services are resolved in the following order:
1. **Route-specific services** (from `UseServices` middleware) - highest priority
2. **App-level services** - fallback if route-specific not found

This allows you to:
- Override app-level services per route
- Provide route-specific service instances
- Compose routes with only the services they need (like Lumora JS!)

### Accessing Services

```go
// Get service (returns error if not found)
userService, err := ctx.Service("userService")
if err != nil {
	return nil, core.NewError(500, "Service unavailable")
}
us := userService.(*UserService)

// Or use MustService (panics if not found) - recommended
us := ctx.MustService("userService").(*UserService)
```

## JSON Body Parsing (useJsonBody)

Similar to Lumora JS's `useJsonBody` hook, you can parse and validate JSON request bodies using the `zog` library:

```go
import (
    z "github.com/Oudwins/zog"
    "github.com/hemant-mann/lumora-go/middleware/usejsonbody"
)

// Define your schema using zog
var userSchema = z.Struct(z.Shape{
    "name":  z.String().Min(3).Max(50),
    "email": z.String().Email(),
    "age":   z.Int().GT(0).LT(150).Optional(),
})

// Define your struct
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   *int   `json:"age,omitempty"`
}

// Use in route
app.Post("/users",
    func(ctx core.Context) (*core.Response, error) {
        // Get parsed and validated body
        user := usejsonbody.GetJsonBody(ctx).(*User)
        
        // Use validated user data
        resp := core.NewResponse().
            WithStatus(201).
            WithBody(map[string]interface{}{
                "message": "User created",
                "user":    user,
            })
        return resp, nil
    },
    usejsonbody.UseJsonBody(userSchema, &User{}),
)

// With custom key
app.Post("/users",
    handler,
    usejsonbody.UseJsonBodyWithKey(userSchema, &User{}, "user"),
)
```

The middleware:
- Parses JSON from request body
- Validates against zog schema
- Returns 400 error response if validation fails
- Stores validated data in context for handler access

**Note**: Handlers return `(*Response, error)` instead of calling `resp.Send(ctx)`. The orchestrator automatically sends the response. If an error is returned, it's handled by the error middleware.

## Header Parsing (useHeaders)

Similar to Lumora JS's `useHeaders` hook, you can parse and validate request headers using the `zog` library:

```go
import (
    z "github.com/Oudwins/zog"
    "github.com/hemant-mann/lumora-go/middleware/useheaders"
)

// Define your schema using zog
// Note: Header names are normalized to lowercase for schema matching
var authHeadersSchema = z.Struct(z.Shape{
    "authorization": z.String().Min(1),
    "x-api-key":     z.String().Optional(),
    "content-type":  z.String().Optional(),
})

// Define your struct
type AuthHeaders struct {
    Authorization string `json:"authorization"`
    APIKey        string `json:"x-api-key,omitempty"`
    ContentType   string `json:"content-type,omitempty"`
}

// Use in route
app.Get("/api/protected",
    func(ctx core.Context) (*core.Response, error) {
        // Get parsed and validated headers
        headers := useheaders.GetHeaders(ctx).(*AuthHeaders)
        
        // Use validated headers
        resp := core.NewResponse().
            WithStatus(200).
            WithBody(map[string]interface{}{
                "message": "Access granted",
                "hasToken": headers.Authorization != "",
            })
        return resp, nil
    },
    useheaders.UseHeaders(authHeadersSchema, &AuthHeaders{}),
)

// With custom key
app.Get("/api/protected",
    handler,
    useheaders.UseHeadersWithKey(authHeadersSchema, &AuthHeaders{}, "authHeaders"),
)

// Combining useHeaders with useJsonBody
app.Post("/api/users",
    func(ctx core.Context) (*core.Response, error) {
        headers := useheaders.GetHeaders(ctx).(*AuthHeaders)
        user := usejsonbody.GetJsonBody(ctx).(*User)
        
        resp := core.NewResponse().
            WithStatus(201).
            WithBody(map[string]interface{}{
                "message": "User created",
                "user":    user,
                "auth":    headers.Authorization != "",
            })
        return resp, nil
    },
    useheaders.UseHeaders(authHeadersSchema, &AuthHeaders{}),
    usejsonbody.UseJsonBody(userSchema, &User{}),
)
```

The middleware:
- Parses headers from request
- Normalizes header names to lowercase for consistent matching
- Validates against zog schema
- Returns 400 error response if validation fails
- Stores validated headers in context for handler access

**Note**: HTTP headers are case-insensitive, so the middleware normalizes them to lowercase. Use lowercase keys in your zog schema (e.g., `"authorization"` not `"Authorization"`).

## License

MIT License with exclusion clause. See [LICENSE](LICENSE) file for details.

**Important**: Cloudstuff Technology Private Limited and its products (including but not limited to Trackier, Affnook, Apptrove, or any products linked with these entities) are expressly prohibited from using this code in any format.

