package core

import (
	"encoding/json"
	"fmt"
)

// Response represents a standardized response structure
type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success creates a success response
func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Data: data,
	}
}

// SuccessWithMessage creates a success response with a message
func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response struct
func ErrorResponse(code int, message string) *Response {
	return &Response{
		Code:  code,
		Error: message,
	}
}

// JSONResponse sends a JSON response through the context
func JSONResponse(ctx Context, code int, data interface{}) error {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	
	encoder := json.NewEncoder(ctx.Response())
	return encoder.Encode(data)
}

// SuccessResponse sends a success response
func SuccessResponse(ctx Context, data interface{}) error {
	return JSONResponse(ctx, 200, Success(data))
}

// SendErrorResponse sends an error response
func SendErrorResponse(ctx Context, code int, message string) error {
	return JSONResponse(ctx, code, ErrorResponse(code, message))
}

// SendResponse sends a response using the Response struct
func SendResponse(ctx Context, resp *Response) error {
	return JSONResponse(ctx, resp.Code, resp)
}

// StringResponse sends a string response
func StringResponse(ctx Context, code int, format string, values ...interface{}) error {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	_, err := fmt.Fprintf(ctx.Response(), format, values...)
	return err
}

