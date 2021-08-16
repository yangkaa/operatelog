package restfulutil

import (
	"github.com/emicklei/go-restful"
)

// ErrReadEntity -
type ErrReadEntity struct {
	err error
}

func (e ErrReadEntity) Error() string {
	return e.err.Error()
}

// Result represents a response for restful api.
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error return error struct
func Error(resp *restful.Response, statusCode int, err error) {
	result := &Result{
		Code: statusCode,
		Msg:  err.Error(),
	}
	resp.WriteHeaderAndJson(statusCode, result, restful.MIME_JSON)
}
