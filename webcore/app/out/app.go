package out

import (
	"runtime/debug"
	"strings"
)

var Environment = "development"

func SetEnvironment(env string) {
	Environment = env
}

// Response represents a standard API response
type Response struct {
	HttpCode   int      `json:"httpCode,omitempty"`
	ErrorCode  int      `json:"errorCode,omitempty"`
	ErrorName  string   `json:"errorName,omitempty"`
	Message    string   `json:"message,omitempty"`
	Data       any      `json:"data,omitempty"`
	StackTrace []string `json:"stack,omitempty"`
	Details    *string  `json:"details,omitempty"`
}

func newResponse(response *Response) *Response {
	if response.HttpCode == 0 {
		response.HttpCode = 200
	}

	if response.ErrorCode != 0 && response.ErrorName != "" && Environment == "development" {
		response.StackTrace = strings.Split(string(debug.Stack()), "\n")
	}

	return response
}

// Error implements the error interface
func (e *Response) Error() string {
	return e.Message
}

// SuccessData creates a success response
func SuccessData(data any) *Response {
	return &Response{
		Data: data,
	}
}

func SuccessMessage(msg string) *Response {
	return &Response{
		Message: msg,
	}
}

// SuccessData creates a success response
func SuccessDataMessage(data any, msg string) *Response {
	return &Response{
		Data:    data,
		Message: msg,
	}
}

func Error(httpCode int, errorCode int, errorName string, message string) *Response {
	return newResponse(&Response{
		HttpCode:  httpCode,
		ErrorCode: errorCode,
		ErrorName: errorName,
		Message:   message,
	})
}

func ErrorDetail(httpCode int, errorCode int, errorName string, message string, detail error) *Response {
	var d *string
	if detail != nil {
		d1 := detail.Error()
		d = &d1
	}

	return newResponse(&Response{
		HttpCode:  httpCode,
		ErrorCode: errorCode,
		ErrorName: errorName,
		Message:   message,
		Details:   d,
	})
}
