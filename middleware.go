package jsonserv

import (
	"log"
	"time"
)

const (
	MaxBodySize = "max_body_size"
	StartTime   = "start_time"
	DebugFlag   = "debug"
)

// instance of middleware

type Middleware interface {
	Ingress(ctx interface{}, req *Request, res *Response)
	Egress(ctx interface{}, req *Request, res *Response)
}

// collection of middleware

type middlewares []Middleware

func (m middlewares) Ingress(ctx interface{}, req *Request, res *Response) {
	for _, middleware := range m {
		middleware.Ingress(ctx, req, res)
	}
}

func (m middlewares) Egress(ctx interface{}, req *Request, res *Response) {
	for i := len(m) - 1; i >= 0; i-- {
		m[i].Egress(ctx, req, res)
	}
}

// built-in middleware

// middleware- LoggingMiddleware

type loggingMiddleware struct {
	logIngress bool
}

func NewLoggingMiddleware(logIngress bool) Middleware {
	return &loggingMiddleware{
		logIngress: logIngress,
	}
}

func (m loggingMiddleware) Ingress(ctx interface{}, req *Request, res *Response) {
	req.SetMiddlewareVar(StartTime, time.Now())
	if m.logIngress {
		log.Printf("← %s %s", req.Method(), req.URL())
	}
}

func (m loggingMiddleware) Egress(ctx interface{}, req *Request, res *Response) {
	start := req.GetMiddlewareVar(StartTime).(time.Time)
	if res.Err != nil {
		log.Printf("→ ERROR %d %s (%s): %v", res.Code, req, time.Now().Sub(start), res.Err)
	} else {
		log.Printf("→ %d %s (%s)", res.Code, req, time.Now().Sub(start))
	}
}

// NewDebugFlagMiddleware creates a middleware that sets a debug flag
// Debug mode will enable error messages in 500 ress
func NewDebugFlagMiddleware(debug bool) Middleware {
	return NewStaticValueMiddleware(DebugFlag, debug)
}

// NewMaxRequestSizeMiddleware creates a middleware that sets the maximum size read of incoming reqs
func NewMaxRequestSizeMiddleware(maxRequestSize int64) Middleware {
	return NewStaticValueMiddleware(MaxBodySize, maxRequestSize)
}

// middleware- StaticValueMiddleware

type staticValueMiddleware struct {
	key   string
	value interface{}
}

func NewStaticValueMiddleware(key string, value interface{}) Middleware {
	return &staticValueMiddleware{
		key:   key,
		value: value,
	}
}

func (m staticValueMiddleware) Ingress(ctx interface{}, req *Request, res *Response) {
	req.SetMiddlewareVar(m.key, m.value)
}

func (m staticValueMiddleware) Egress(ctx interface{}, req *Request, res *Response) {
}

// middleware- FactoryValueMiddleware

type Factory func(ctx interface{}, req *Request, res *Response) interface{}

type dynamicValueMiddleware struct {
	key     string
	factory Factory
}

func NewDynamicValueMiddleware(key string, factory Factory) Middleware {
	return &dynamicValueMiddleware{
		key:     key,
		factory: factory,
	}
}

func (m dynamicValueMiddleware) Ingress(ctx interface{}, req *Request, res *Response) {
	req.SetMiddlewareVar(m.key, m.factory(ctx, req, res))
}

func (m dynamicValueMiddleware) Egress(ctx interface{}, req *Request, res *Response) {
}
