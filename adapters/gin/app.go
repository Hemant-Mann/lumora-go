package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
)

type App struct {
	engine      *gin.Engine
	middlewares []core.Middleware
	services    *services.Container
}

// New creates a new gin adapter app
func New() *App {
	return &App{
		engine:      gin.New(),
		middlewares: []core.Middleware{},
		services:    services.NewContainer(),
	}
}

func (a *App) Use(middleware ...core.Middleware) {
	a.middlewares = append(a.middlewares, middleware...)
}

func (a *App) Handle(method, path string, handler core.Handler, middlewares ...core.Middleware) {
	// Combine app-level and route-level middlewares
	allMiddlewares := append(a.middlewares, middlewares...)
	
	// Apply middlewares to handler
	finalHandler := core.Apply(handler, allMiddlewares...)
	
		// Convert to gin handler
		ginHandler := func(ginCtx *gin.Context) {
			ctx := NewContext(ginCtx, a.services)
			if err := finalHandler(ctx); err != nil {
				// Error will be handled by error middleware if present
				ginCtx.Error(err)
			}
		}
	
	// Register with gin
	switch method {
	case "GET":
		a.engine.GET(path, ginHandler)
	case "POST":
		a.engine.POST(path, ginHandler)
	case "PUT":
		a.engine.PUT(path, ginHandler)
	case "DELETE":
		a.engine.DELETE(path, ginHandler)
	case "PATCH":
		a.engine.PATCH(path, ginHandler)
	default:
		a.engine.Handle(method, path, ginHandler)
	}
}

func (a *App) Get(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle("GET", path, handler, middlewares...)
}

func (a *App) Post(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle("POST", path, handler, middlewares...)
}

func (a *App) Put(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle("PUT", path, handler, middlewares...)
}

func (a *App) Delete(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle("DELETE", path, handler, middlewares...)
}

func (a *App) Patch(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle("PATCH", path, handler, middlewares...)
}

func (a *App) Services() *services.Container {
	return a.services
}

func (a *App) Start(addr string) error {
	return a.engine.Run(addr)
}

