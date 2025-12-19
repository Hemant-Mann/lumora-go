package nethttp

import (
	"fmt"
	"net/http"

	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
)

type App struct {
	mux         *http.ServeMux
	middlewares []core.Middleware
	router      *Router
	services    *services.Container
}

// New creates a new net/http adapter app
func New() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []core.Middleware{},
		router:      NewRouter(),
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
	
	// Register with router
	a.router.Handle(method, path, func(ctx core.Context) error {
		return finalHandler(ctx)
	})
}

func (a *App) Get(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle(http.MethodGet, path, handler, middlewares...)
}

func (a *App) Post(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle(http.MethodPost, path, handler, middlewares...)
}

func (a *App) Put(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle(http.MethodPut, path, handler, middlewares...)
}

func (a *App) Delete(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle(http.MethodDelete, path, handler, middlewares...)
}

func (a *App) Patch(path string, handler core.Handler, middlewares ...core.Middleware) {
	a.Handle(http.MethodPatch, path, handler, middlewares...)
}

func (a *App) Services() *services.Container {
	return a.services
}

func (a *App) Start(addr string) error {
	// Register router handler with mux
	a.mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req, res, a.services)
		
		// Set app-level services in context for UseServices middleware
		ctx.Set("_app_services", a.services)
		
		// Try to match route
		handler, params := a.router.Match(req.Method, req.URL.Path)
		if handler == nil {
			http.NotFound(res, req)
			return
		}
		
		// Set path parameters
		if ctxImpl, ok := ctx.(*contextImpl); ok {
			ctxImpl.SetParams(params)
		}
		
		// Execute handler
		if err := handler(ctx); err != nil {
			// Error handling will be done by error middleware if present
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	})
	
	fmt.Printf("Server starting on %s\n", addr)
	return http.ListenAndServe(addr, a.mux)
}

