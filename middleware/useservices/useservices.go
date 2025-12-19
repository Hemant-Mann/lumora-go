package useservices

import (
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
)

// UseServices creates a middleware that injects route-specific services into the context
// These services are available only for this route and override app-level services with the same name
func UseServices(routeServices map[string]interface{}) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx core.Context) error {
			// Create a scoped service container for this route
			scopedContainer := services.NewContainer()

			// Copy app-level services first (if available)
			if appServices := getAppServices(ctx); appServices != nil {
				for name, service := range appServices.All() {
					scopedContainer.Register(name, service)
				}
			}

			// Override with route-specific services
			for name, service := range routeServices {
				scopedContainer.Register(name, service)
			}

			// Set the scoped container in context
			setScopedServices(ctx, scopedContainer)

			return next(ctx)
		}
	}
}

// UseService is a convenience function for a single service
func UseService(name string, service interface{}) core.Middleware {
	return UseServices(map[string]interface{}{name: service})
}

// getAppServices extracts the app-level service container from context
// This is a helper that works with the context's internal structure
func getAppServices(ctx core.Context) *services.Container {
	// Try to get from context values (set by adapters)
	if container, ok := ctx.Get("_app_services"); ok {
		if svcs, ok := container.(*services.Container); ok {
			return svcs
		}
	}
	return nil
}

// setScopedServices sets the scoped service container in context
func setScopedServices(ctx core.Context, container *services.Container) {
	ctx.Set("_scoped_services", container)
}

// GetScopedServices retrieves the scoped service container from context
func GetScopedServices(ctx core.Context) *services.Container {
	if container, ok := ctx.Get("_scoped_services"); ok {
		if svcs, ok := container.(*services.Container); ok {
			return svcs
		}
	}
	return nil
}
