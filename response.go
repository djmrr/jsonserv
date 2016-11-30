package jsonserv

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
	Close()
}

type ResponseWriterCloser struct {
	http.ResponseWriter
}

func (r ResponseWriterCloser) Close(){}

type Response struct {
	Code    int
	Err     error
	Body    interface{}
	Writer  ResponseWriter
}

func newWrappedResponse(w http.ResponseWriter) *Response {
	return newResponse(ResponseWriterCloser{w})
}

func newResponse(w ResponseWriter) *Response {
	return &Response{
		Code:   http.StatusOK,
		Writer: w,
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
	r.Writer.Header().Add(key, value)
	return r
}
