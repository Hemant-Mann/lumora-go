package nethttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
)

type contextImpl struct {
	req        *http.Request
	res        http.ResponseWriter
	params     map[string]string
	queryCache map[string]string
	values     map[string]interface{}
	statusCode int
	services   *services.Container
}

// NewContext creates a new context from http.Request and http.ResponseWriter
func NewContext(req *http.Request, res http.ResponseWriter, svcs *services.Container) core.Context {
	return &contextImpl{
		req:        req,
		res:        res,
		params:     make(map[string]string),
		queryCache: make(map[string]string),
		values:     make(map[string]interface{}),
		statusCode: 200,
		services:   svcs,
	}
}

func (c *contextImpl) Request() *http.Request {
	return c.req
}

func (c *contextImpl) Response() http.ResponseWriter {
	return c.res
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
	if val, ok := c.queryCache[name]; ok {
		return val
	}
	val := c.req.URL.Query().Get(name)
	c.queryCache[name] = val
	return val
}

func (c *contextImpl) Header(name string) string {
	return c.req.Header.Get(name)
}

func (c *contextImpl) SetHeader(name, value string) {
	c.res.Header().Set(name, value)
}

func (c *contextImpl) Status(code int) {
	c.statusCode = code
	c.res.WriteHeader(code)
}

func (c *contextImpl) JSON(code int, data interface{}) error {
	// Set header BEFORE calling Status() (which calls WriteHeader())
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.res)
	return encoder.Encode(data)
}

func (c *contextImpl) String(code int, format string, values ...interface{}) error {
	// Set header BEFORE calling Status() (which calls WriteHeader())
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := fmt.Fprintf(c.res, format, values...)
	return err
}

func (c *contextImpl) BindJSON(dest interface{}) error {
	decoder := json.NewDecoder(c.req.Body)
	return decoder.Decode(dest)
}

func (c *contextImpl) Context() context.Context {
	return c.req.Context()
}

func (c *contextImpl) WithContext(ctx context.Context) core.Context {
	newReq := c.req.WithContext(ctx)
	return &contextImpl{
		req:        newReq,
		res:        c.res,
		params:     c.params,
		queryCache: c.queryCache,
		values:     c.values,
		statusCode: c.statusCode,
		services:   c.services,
	}
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
	// Read the request body
	body, err := io.ReadAll(c.req.Body)
	if err != nil {
		return nil, err
	}
	// Restore the body so it can be read again if needed
	c.req.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}
