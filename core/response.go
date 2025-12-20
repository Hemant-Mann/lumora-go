package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Cookie represents an HTTP cookie
type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// Response represents a flexible response structure
type Response struct {
	StatusCode int
	Headers    map[string]string
	Cookies    []Cookie
	Body       interface{}
}

// NewResponse creates a new response with default values
func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
		Headers:    make(map[string]string),
		Cookies:    []Cookie{},
		Body:       nil,
	}
}

// WithStatus sets the status code
func (r *Response) WithStatus(code int) *Response {
	r.StatusCode = code
	return r
}

// WithHeader sets a header
func (r *Response) WithHeader(name, value string) *Response {
	r.Headers[name] = value
	return r
}

// WithCookie adds a cookie
func (r *Response) WithCookie(cookie Cookie) *Response {
	r.Cookies = append(r.Cookies, cookie)
	return r
}

// WithBody sets the response body
func (r *Response) WithBody(body interface{}) *Response {
	r.Body = body
	return r
}

// Send sends the response through the context
func (r *Response) Send(ctx Context) error {
	// IMPORTANT: Set headers and cookies BEFORE calling Status()
	// In net/http, once WriteHeader() is called, headers cannot be modified

	// Set headers first
	for name, value := range r.Headers {
		ctx.SetHeader(name, value)
	}

	// Set cookies
	for _, cookie := range r.Cookies {
		httpCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			MaxAge:   cookie.MaxAge,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: cookie.SameSite,
		}
		http.SetCookie(ctx.Response(), httpCookie)
	}

	// Handle body
	if r.Body == nil {
		ctx.Status(r.StatusCode)
		return nil
	}

	// Check if body is a string - send as plain text
	if bodyStr, ok := r.Body.(string); ok {
		// Set content type if not already set
		if _, exists := r.Headers["Content-Type"]; !exists {
			ctx.SetHeader("Content-Type", "text/plain")
		}
		ctx.Status(r.StatusCode)
		_, err := fmt.Fprint(ctx.Response(), bodyStr)
		return err
	}

	// Otherwise, send as JSON
	// Set content type if not already set
	if _, exists := r.Headers["Content-Type"]; !exists {
		ctx.SetHeader("Content-Type", "application/json")
	}
	ctx.Status(r.StatusCode)

	encoder := json.NewEncoder(ctx.Response())
	return encoder.Encode(r.Body)
}

// SendResponse is a helper to send a Response struct
func SendResponse(ctx Context, resp *Response) error {
	return resp.Send(ctx)
}

// HandleResponse is the orchestrator function that handles response and error
// If error is not nil, it returns the error (to be handled by error middleware)
// If error is nil and response is not nil, it sends the response
func HandleResponse(ctx Context, resp *Response, err error) error {
	if err != nil {
		return err
	}
	if resp != nil {
		return resp.Send(ctx)
	}
	return nil
}

// Helper functions for common response patterns

// JSON sends a JSON response
func JSON(ctx Context, code int, data interface{}) error {
	resp := NewResponse().
		WithStatus(code).
		WithHeader("Content-Type", "application/json").
		WithBody(data)
	return resp.Send(ctx)
}

// Text sends a plain text response
func Text(ctx Context, code int, text string) error {
	resp := NewResponse().
		WithStatus(code).
		WithHeader("Content-Type", "text/plain").
		WithBody(text)
	return resp.Send(ctx)
}

// String sends a formatted string response (for backward compatibility)
func String(ctx Context, code int, format string, values ...interface{}) error {
	text := fmt.Sprintf(format, values...)
	return Text(ctx, code, text)
}
