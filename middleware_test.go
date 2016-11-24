package jsonserv

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

type countingmiddleware struct {
	ingress, egress int
}

func (m *countingmiddleware) Ingress(ctx interface{}, req *Request, res *Response) {
	m.ingress++
}
func (m *countingmiddleware) Egress(ctx interface{}, req *Request, res *Response) {
	m.egress++
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestMiddlewares_Ingress(t *testing.T) {
	c := new(countingmiddleware)
	m := middlewares{c}

	m.Ingress(nil, nil, nil)
	if c.ingress != 1 {
		t.Fatalf("Ingress incorrect: %d", c.ingress)
	}
}

func TestMiddlewares_Egress(t *testing.T) {
	c := new(countingmiddleware)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Egress(nil, req, res)
	if c.egress != 1 {
		t.Fatalf("Egress incorrect: %d", c.egress)
	}
}

func TestLoggingMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewLoggingMiddleware(true)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar(StartTime).(time.Time).IsZero() {
		t.Fatal("Unexpected start time")
	}
	m.Egress(nil, req, res)
}

func TestLoggingMiddleware_With_Err(t *testing.T) {
	c := NewLoggingMiddleware(true)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	res.Err = errors.New("should be printed")
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar(StartTime).(time.Time).IsZero() {
		t.Fatal("Unexpected start time")
	}
	m.Egress(nil, req, res)
}

func TestDebugFlagMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewDebugFlagMiddleware(true)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Ingress(nil, req, res)
	if !req.GetMiddlewareVar(DebugFlag).(bool) {
		t.Fatal("Unexpected debug")
	}
	m.Egress(nil, req, res)
}

func TestMaxRequestSizeMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewMaxRequestSizeMiddleware(5000)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar(MaxBodySize).(int64) != 5000 {
		t.Fatal("Unexpected max body size")
	}
	m.Egress(nil, req, res)
}

func TestStaticValueMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewStaticValueMiddleware("foo", "bar")
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar("foo").(string) != "bar" {
		t.Fatal("Unexpected foo")
	}
	m.Egress(nil, req, res)
}

func TestDynamicValueMiddleware_RequestVar_Is_Set(t *testing.T) {
	i := 0
	c := NewDynamicValueMiddleware("foo", func(ctx interface{}, req *Request, res *Response) interface{} {
		i++
		return i
	})
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse()
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar("foo").(int) != 1 {
		t.Fatal("Unexpected foo")
	}
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar("foo").(int) != 2 {
		t.Fatal("Unexpected foo")
	}
	m.Egress(nil, req, res)
}
