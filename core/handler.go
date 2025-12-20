package core

// Handler is the main handler function type
// Returns a Response and an error. The orchestrator will handle sending the response.
type Handler func(Context) (*Response, error)

// Middleware is a function that wraps a handler
type Middleware func(Handler) Handler

// Compose chains multiple middlewares together
func Compose(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Apply applies middleware to a handler
func Apply(handler Handler, middlewares ...Middleware) Handler {
	return Compose(middlewares...)(handler)
}

