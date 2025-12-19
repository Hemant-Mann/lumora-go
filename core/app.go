package core

// App represents the main application interface
type App interface {
	// Use adds middleware to the application
	Use(middleware ...Middleware)
	
	// Handle registers a handler for a specific method and path
	Handle(method, path string, handler Handler, middlewares ...Middleware)
	
	// Get registers a GET route
	Get(path string, handler Handler, middlewares ...Middleware)
	
	// Post registers a POST route
	Post(path string, handler Handler, middlewares ...Middleware)
	
	// Put registers a PUT route
	Put(path string, handler Handler, middlewares ...Middleware)
	
	// Delete registers a DELETE route
	Delete(path string, handler Handler, middlewares ...Middleware)
	
	// Patch registers a PATCH route
	Patch(path string, handler Handler, middlewares ...Middleware)
	
	// Start starts the server
	Start(addr string) error
}

