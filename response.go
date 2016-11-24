package jsonserv

import "net/http"

type Response struct {
	Code    int
	Err     error
	Body    interface{}
	headers map[string]string
}

func newResponse() *Response {
	return &Response{
		Code: http.StatusOK,
	}
}

func (r *Response) Done(code int, body interface{}) *Response {
	r.Code = code
	r.Body = body
	return r
}

func (r *Response) Empty(code int) *Response {
	return r.Done(code, nil)
}

func (r *Response) Ok(body interface{}) *Response {
	return r.Done(http.StatusOK, body)
}

func (r *Response) Error(err error) *Response {
	r.Code = http.StatusInternalServerError
	r.Err = err
	r.Body = nil
	return r
}

func (r *Response) AddHeader(key, value string) *Response {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
	return r
}
