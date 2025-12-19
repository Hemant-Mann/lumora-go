package fasthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
	"github.com/valyala/fasthttp"
)

type contextImpl struct {
	ctx      *fasthttp.RequestCtx
	params   map[string]string
	values   map[string]interface{}
	reqCtx   context.Context
	services *services.Container
}

// NewContext creates a new context from fasthttp.RequestCtx
func NewContext(ctx *fasthttp.RequestCtx, svcs *services.Container) core.Context {
	return &contextImpl{
		ctx:      ctx,
		params:   make(map[string]string),
		values:   make(map[string]interface{}),
		reqCtx:   context.Background(),
		services: svcs,
	}
}

func (c *contextImpl) Request() *http.Request {
	// Convert fasthttp request to net/http request
	req := new(http.Request)
	req.Method = string(c.ctx.Method())
	req.URL = &url.URL{
		Path:     string(c.ctx.Path()),
		RawQuery: string(c.ctx.QueryArgs().QueryString()),
	}
	req.Header = make(http.Header)
	c.ctx.Request.Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})
	return req
}

func (c *contextImpl) Response() http.ResponseWriter {
	// Return a wrapper that writes to fasthttp response
	return &responseWriter{ctx: c.ctx}
}

func (c *contextImpl) Get(key string) (interface{}, bool) {
	val, ok := c.values[key]
	return val, ok
}

func (c *contextImpl) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *contextImpl) Param(name string) string {
	return c.params[name]
}

func (c *contextImpl) Query(name string) string {
	return string(c.ctx.QueryArgs().Peek(name))
}

func (c *contextImpl) Header(name string) string {
	return string(c.ctx.Request.Header.Peek(name))
}

func (c *contextImpl) SetHeader(name, value string) {
	c.ctx.Response.Header.Set(name, value)
}

func (c *contextImpl) Status(code int) {
	c.ctx.SetStatusCode(code)
}

func (c *contextImpl) JSON(code int, data interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.ctx.Response.BodyWriter())
	return encoder.Encode(data)
}

func (c *contextImpl) String(code int, format string, values ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := fmt.Fprintf(c.ctx, format, values...)
	return err
}

func (c *contextImpl) BindJSON(dest interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(c.ctx.PostBody()))
	return decoder.Decode(dest)
}

func (c *contextImpl) Context() context.Context {
	return c.reqCtx
}

func (c *contextImpl) WithContext(ctx context.Context) core.Context {
	newCtx := &contextImpl{
		ctx:      c.ctx,
		params:   c.params,
		values:   c.values,
		reqCtx:   ctx,
		services: c.services,
	}
	return newCtx
}

func (c *contextImpl) Service(name string) (interface{}, error) {
	// Check for scoped services first (route-specific)
	if scopedContainer, ok := c.values["_scoped_services"]; ok {
		if svcs, ok := scopedContainer.(*services.Container); ok {
			if service, err := svcs.Get(name); err == nil {
				return service, nil
			}
		}
	}

	// Fall back to app-level services
	if c.services == nil {
		return nil, fmt.Errorf("services container not available")
	}
	return c.services.Get(name)
}

func (c *contextImpl) MustService(name string) interface{} {
	// Check for scoped services first (route-specific)
	if scopedContainer, ok := c.values["_scoped_services"]; ok {
		if svcs, ok := scopedContainer.(*services.Container); ok {
			if svcs.Has(name) {
				return svcs.MustGet(name)
			}
		}
	}

	// Fall back to app-level services
	if c.services == nil {
		panic("services container not available")
	}
	return c.services.MustGet(name)
}

// SetParams sets the path parameters (used by router)
func (c *contextImpl) SetParams(params map[string]string) {
	c.params = params
}

func (c *contextImpl) RequestBody() ([]byte, error) {
	// In fasthttp, PostBody() returns the body as bytes
	return c.ctx.PostBody(), nil
}

// responseWriter wraps fasthttp.RequestCtx to implement http.ResponseWriter
type responseWriter struct {
	ctx *fasthttp.RequestCtx
}

func (w *responseWriter) Header() http.Header {
	header := make(http.Header)
	w.ctx.Response.Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.ctx.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ctx.SetStatusCode(statusCode)
}
