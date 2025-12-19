package gin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hemant-mann/lumora-go/core"
	"github.com/hemant-mann/lumora-go/services"
)

type contextImpl struct {
	ctx      *gin.Context
	services *services.Container
}

// NewContext creates a new context from gin.Context
func NewContext(ctx *gin.Context, svcs *services.Container) core.Context {
	return &contextImpl{ctx: ctx, services: svcs}
}

func (c *contextImpl) Request() *http.Request {
	return c.ctx.Request
}

func (c *contextImpl) Response() http.ResponseWriter {
	return c.ctx.Writer
}

func (c *contextImpl) Get(key string) (interface{}, bool) {
	val, ok := c.ctx.Get(key)
	return val, ok
}

func (c *contextImpl) Set(key string, value interface{}) {
	c.ctx.Set(key, value)
}

func (c *contextImpl) Param(name string) string {
	return c.ctx.Param(name)
}

func (c *contextImpl) Query(name string) string {
	return c.ctx.Query(name)
}

func (c *contextImpl) Header(name string) string {
	return c.ctx.GetHeader(name)
}

func (c *contextImpl) SetHeader(name, value string) {
	c.ctx.Header(name, value)
}

func (c *contextImpl) Status(code int) {
	c.ctx.Status(code)
}

func (c *contextImpl) JSON(code int, data interface{}) error {
	c.ctx.JSON(code, data)
	return nil
}

func (c *contextImpl) String(code int, format string, values ...interface{}) error {
	c.ctx.String(code, format, values...)
	return nil
}

func (c *contextImpl) BindJSON(dest interface{}) error {
	return c.ctx.ShouldBindJSON(dest)
}

func (c *contextImpl) Context() context.Context {
	return c.ctx.Request.Context()
}

func (c *contextImpl) WithContext(ctx context.Context) core.Context {
	newGinCtx := c.ctx.Copy()
	newGinCtx.Request = newGinCtx.Request.WithContext(ctx)
	return &contextImpl{ctx: newGinCtx, services: c.services}
}

func (c *contextImpl) Service(name string) (interface{}, error) {
	if c.services == nil {
		return nil, fmt.Errorf("services container not available")
	}
	return c.services.Get(name)
}

func (c *contextImpl) MustService(name string) interface{} {
	if c.services == nil {
		panic("services container not available")
	}
	return c.services.MustGet(name)
}

