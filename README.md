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
	
	// Define routes
	app.Get("/", func(ctx core.Context) error {
		return core.SuccessResponse(ctx, map[string]string{
			"message": "Hello, World!",
		})
	})
	
	app.Get("/users/:id", func(ctx core.Context) error {
		id := ctx.Param("id")
		return core.SuccessResponse(ctx, map[string]string{
			"id": id,
		})
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
		return core.SuccessResponse(ctx, map[string]string{
			"message": "Hello from Gin!",
		})
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
		return core.SuccessResponse(ctx, map[string]string{
			"message": "Hello from FastHTTP!",
		})
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

## Response Helpers

```go
// Success response
core.SuccessResponse(ctx, data)

// Error response
core.ErrorResponse(ctx, 400, "Bad Request")

// Custom JSON response
core.JSONResponse(ctx, 201, customData)
```

## License

MIT

