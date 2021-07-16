package restfulutil

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
