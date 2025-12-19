package fasthttp

import (
	"fmt"

	"github.com/hemant-mann/lumora-go/core"
	"github.com/valyala/fasthttp"
)

type App struct {
	server      *fasthttp.Server
	middlewares []core.Middleware
	router      *Router
}

// New creates a new fasthttp adapter app
func New() *App {
	app := &App{
		middlewares: []core.Middleware{},
		router:      NewRouter(),
	}

	app.server = &fasthttp.Server{
		Handler: app.handleRequest,
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

	// Register with router
	a.router.Handle(method, path, func(ctx core.Context) error {
		return finalHandler(ctx)
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

func (a *App) Start(addr string) error {
	fmt.Printf("Server starting on %s\n", addr)
	return a.server.ListenAndServe(addr)
}

func (a *App) handleRequest(ctx *fasthttp.RequestCtx) {
	coreCtx := NewContext(ctx)

	// Try to match route
	handler, params := a.router.Match(string(ctx.Method()), string(ctx.Path()))
	if handler == nil {
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}

	// Set path parameters
	if ctxImpl, ok := coreCtx.(*contextImpl); ok {
		ctxImpl.SetParams(params)
	}

	// Execute handler
	if err := handler(coreCtx); err != nil {
		// Error handling will be done by error middleware if present
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
