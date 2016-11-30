package jsonserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Request struct {
	raw *http.Request
	ctx map[string]interface{}
}

func newRequest(r *http.Request) *Request {
	req := &Request{
		raw: r,
	}
	return req
}

func (r *Request) Method() string {
	return r.raw.Method
}

func (r *Request) URL() *url.URL {
	return r.raw.URL
}

func (r *Request) Header() http.Header {
	return r.raw.Header
}

func (r *Request) GetMiddlewareVar(key string) interface{} {
	if r.ctx == nil {
		return nil
	}
	return r.ctx[key]
}

func (r *Request) GetOptionalMiddlewareVar(key string, fallback interface{}) interface{} {
	if r.ctx == nil {
		return fallback
	}
	val, ok := r.ctx[key]
	if !ok {
		return fallback
	}
	return val
}

func (r *Request) SetMiddlewareVar(key string, value interface{}) {
	if r.ctx == nil {
		r.ctx = make(map[string]interface{})
	}
	r.ctx[key] = value
}

func (r Request) String() string {
	return fmt.Sprintf("%s %s", r.raw.Method, r.raw.URL)
}

func (r *Request) GetPathVars() map[string]string {
	return mux.Vars(r.raw)
}

func (r *Request) GetPathVar(key string, fallback string) string {
	vars := mux.Vars(r.raw)
	if vars == nil {
		return fallback
	}
	value, ok := vars[key]
	if !ok {
		return fallback
	} else {
		return value
	}
}

func (r *Request) ParseBody(v interface{}) error {
	maxRequestSize := r.GetOptionalMiddlewareVar(MaxBodySize, int64(0)).(int64)
	var reader io.Reader
	if maxRequestSize == 0 {
		reader = r.raw.Body
	} else if r.raw.ContentLength > maxRequestSize {
		return errors.New("Request body too large")
	} else {
		reader = io.LimitReader(r.raw.Body, maxRequestSize)
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if err := r.raw.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
