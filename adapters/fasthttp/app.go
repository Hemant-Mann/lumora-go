package fasthttp

import (
	"fmt"

	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
	"github.com/valyala/fasthttp"
)

type App struct {
	server      *fasthttp.Server
	middlewares []core.Middleware
	router      *Router
	services    *services.Container
}

// New creates a new fasthttp adapter app
func New() *App {
	app := &App{
		middlewares: []core.Middleware{},
		router:      NewRouter(),
		services:    services.NewContainer(),
	}

	// Wrap the router handler with our middleware handler
	app.server = &fasthttp.Server{
		Handler: app.wrapHandler(app.router.Handler()),
	}

	return app
}

func (a *App) Use(middleware ...core.Middleware) {
	a.middlewares = append(a.middlewares, middleware...)
}

func (a *App) Handle(method, path string, handler core.Handler, middlewares ...core.Middleware) {
	// Combine app-level and route-level middlewares
	allMiddlewares := append(a.middlewares, middlewares...)

	// Apply middlewares to handler
	finalHandler := core.Apply(handler, allMiddlewares...)

	// Register with router - create a fasthttp handler that converts context
	a.router.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		coreCtx := NewContext(ctx, a.services)

		// Set app-level services in context for UseServices middleware
		coreCtx.Set("_app_services", a.services)

		// Extract path parameters from UserValues (fasthttp/router stores them here)
		if ctxImpl, ok := coreCtx.(*contextImpl); ok {
			params := make(map[string]string)
			ctx.VisitUserValues(func(key []byte, value any) {
				if str, ok := value.(string); ok {
					params[string(key)] = str
				}
			})
			ctxImpl.SetParams(params)
		}

		// Call our core handler - orchestrator handles response and error
		resp, handlerErr := finalHandler(coreCtx)
		if err := core.HandleResponse(coreCtx, resp, handlerErr); err != nil {
			// If error middleware didn't handle it, send a default error response
			// This should rarely happen if error middleware is properly configured
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
	})
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
	fmt.Printf("Server starting on %s\n", addr)
	return a.server.ListenAndServe(addr)
}

// wrapHandler wraps the router handler to handle errors from our core handlers
func (a *App) wrapHandler(routerHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Call the router handler
		// The router will call our registered handlers which return errors
		// If an error is returned, it will be handled by error middleware
		routerHandler(ctx)
	}
}
