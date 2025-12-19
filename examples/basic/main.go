package main

import (
	"github.com/hemant-mann/lumora-go/adapters/nethttp"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/middleware/cors"
	"github.com/hemant-mann/lumora-go/middleware/errorhandler"
	"github.com/hemant-mann/lumora-go/middleware/logging"
)

// Example service
type UserService struct {
	users map[string]string
}

func NewUserService() *UserService {
	return &UserService{
		users: map[string]string{
			"1": "Alice",
			"2": "Bob",
			"3": "Charlie",
		},
	}
}

func (s *UserService) GetUser(id string) (string, bool) {
	name, ok := s.users[id]
	return name, ok
}

func main() {
	app := nethttp.New()

	// Register services
	userService := NewUserService()
	app.Services().Register("userService", userService)

	// Add global middleware
	app.Use(
		cors.New(cors.DefaultOptions()),
		logging.Simple(),
		errorhandler.Simple(),
	)

	// Define routes
	app.Get("/", func(ctx core.Context) error {
		// Using new Response system - user controls the structure
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]interface{}{
				"message": "Hello, World!",
				"version": "1.0.0",
			})
		return resp.Send(ctx)
	})

	app.Get("/users/:id", func(ctx core.Context) error {
		id := ctx.Param("id")

		// Access service from context
		userService, err := ctx.Service("userService")
		if err != nil {
			return core.NewError(500, "Service not available")
		}

		// Type assert to use the service
		us := userService.(*UserService)
		name, exists := us.GetUser(id)

		if !exists {
			resp := core.NewResponse().
				WithStatus(404).
				WithBody(map[string]string{"error": "User not found"})
			return resp.Send(ctx)
		}

		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]interface{}{
				"id":   id,
				"name": name,
			})
		return resp.Send(ctx)
	})

	app.Post("/users", func(ctx core.Context) error {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := ctx.BindJSON(&user); err != nil {
			resp := core.NewResponse().
				WithStatus(400).
				WithBody(map[string]string{"error": "Invalid JSON"})
			return resp.Send(ctx)
		}

		resp := core.NewResponse().
			WithStatus(201).
			WithBody(map[string]interface{}{
				"message": "User created",
				"user":    user,
			})
		return resp.Send(ctx)
	})

	// Example with plain text response
	app.Get("/text", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithHeader("Content-Type", "text/plain").
			WithBody("This is a plain text response")
		return resp.Send(ctx)
	})

	// Example with cookies
	app.Get("/cookie", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithCookie(core.Cookie{
				Name:     "session",
				Value:    "abc123",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			}).
			WithBody(map[string]string{"message": "Cookie set"})
		return resp.Send(ctx)
	})

	app.Get("/error", func(ctx core.Context) error {
		return core.NewError(500, "This is an error example")
	})

	app.Start(":8080")
}
