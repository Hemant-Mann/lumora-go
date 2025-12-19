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
	app.Get("/", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello, World!",
			})
		return resp.Send(ctx)
	})
	
	app.Get("/users/:id", func(ctx core.Context) error {
		id := ctx.Param("id")
		
		// Access service
		userService := ctx.MustService("userService").(*UserService)
		
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"id": id,
			})
		return resp.Send(ctx)
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
	
	app.Get("/", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello from Gin!",
			})
		return resp.Send(ctx)
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
	
	app.Get("/", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]string{
				"message": "Hello from FastHTTP!",
			})
		return resp.Send(ctx)
	})
	
	app.Start(":8080")
}
```

## Architecture

### Core Package

The `core` package defines framework-agnostic interfaces:
- `Context`: Request/response context abstraction
- `Handler`: Handler function type
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

// In your handler
if someCondition {
	return core.NewError(400, "Bad Request")
}
```

## Response System

The response system gives you full control over your response structure:

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
return resp.Send(ctx)
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

Services provide dependency injection capabilities:

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

// Register services
app.Services().Register("userService", NewUserService())

// Or register by type
app.Services().RegisterByType(NewUserService())

// Access services in handlers
app.Get("/users/:id", func(ctx core.Context) error {
	// Get service (returns error if not found)
	userService, err := ctx.Service("userService")
	if err != nil {
		return core.NewError(500, "Service unavailable")
	}
	
	// Type assert and use
	us := userService.(*UserService)
	// ... use the service
	
	// Or use MustService (panics if not found)
	us := ctx.MustService("userService").(*UserService)
	
	return resp.Send(ctx)
})
```

## License

MIT

