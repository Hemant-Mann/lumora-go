package main

import (
	z "github.com/Oudwins/zog"
	"github.com/hemant-mann/lumora-go/adapters/nethttp"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/middleware/cors"
	"github.com/hemant-mann/lumora-go/middleware/errorhandler"
	"github.com/hemant-mann/lumora-go/middleware/logging"
	"github.com/hemant-mann/lumora-go/middleware/usejsonbody"
	"github.com/hemant-mann/lumora-go/middleware/useservices"
)

// Example services
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

type AuthService struct {
	tokens map[string]bool
}

func NewAuthService() *AuthService {
	return &AuthService{
		tokens: map[string]bool{
			"token123": true,
			"token456": true,
		},
	}
}

func (s *AuthService) ValidateToken(token string) bool {
	return s.tokens[token]
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	app := nethttp.New()

	// Register app-level services (available to all routes)
	app.Services().Register("authService", NewAuthService())

	// Add global middleware
	app.Use(
		cors.New(cors.DefaultOptions()),
		logging.Simple(),
		errorhandler.Simple(),
	)

	// Define routes with route-specific services (like Lumora JS!)
	app.Get("/", func(ctx core.Context) error {
		resp := core.NewResponse().
			WithStatus(200).
			WithBody(map[string]interface{}{
				"message": "Hello, World!",
				"version": "1.0.0",
			})
		return resp.Send(ctx)
	})

	// Route with route-specific services - similar to Lumora JS useServices
	app.Get("/users/:id",
		func(ctx core.Context) error {
			id := ctx.Param("id")

			// Access route-specific service
			userService := ctx.MustService("userService").(*UserService)
			name, exists := userService.GetUser(id)

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
		},
		useservices.UseServices(map[string]interface{}{
			"userService": NewUserService(),
		}),
	)

	// Route with multiple route-specific services
	app.Get("/protected/:id",
		func(ctx core.Context) error {
			// Access route-specific services
			userService := ctx.MustService("userService").(*UserService)
			authService := ctx.MustService("authService").(*AuthService)

			// Check auth token
			token := ctx.Header("Authorization")
			if !authService.ValidateToken(token) {
				resp := core.NewResponse().
					WithStatus(401).
					WithBody(map[string]string{"error": "Unauthorized"})
				return resp.Send(ctx)
			}

			id := ctx.Param("id")
			name, exists := userService.GetUser(id)
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
					"auth": "validated",
				})
			return resp.Send(ctx)
		},
		useservices.UseServices(map[string]interface{}{
			"userService": NewUserService(),
			"authService": NewAuthService(), // Overrides app-level authService for this route
		}),
	)

	var userSchema = z.Struct(z.Shape{
		"name":  z.String().Min(3).Max(50),
		"email": z.String().Email(),
		"age":   z.Int().GT(0).LT(150).Optional(),
	})

	// Using UseService helper for single service
	app.Post("/users",
		func(ctx core.Context) error {
			user := usejsonbody.GetJsonBody(ctx).(*User)
			resp := core.NewResponse().
				WithStatus(201).
				WithBody(map[string]any{
					"message": "User created",
					"user":    user,
				})
			return resp.Send(ctx)
		},
		useservices.UseService("userService", NewUserService()),
		usejsonbody.UseJsonBody(userSchema, &User{}),
	)

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
