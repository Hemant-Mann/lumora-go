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
		return core.SuccessResponse(ctx, map[string]interface{}{
			"id":   id,
			"name": "User " + id,
		})
	})
	
	app.Post("/users", func(ctx core.Context) error {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		
		if err := ctx.BindJSON(&user); err != nil {
			return core.NewError(400, "Invalid JSON")
		}
		
		return core.SuccessResponse(ctx, map[string]interface{}{
			"message": "User created",
			"user":    user,
		})
	})
	
	app.Get("/error", func(ctx core.Context) error {
		return core.NewError(500, "This is an error example")
	})
	
	app.Start(":8080")
}

