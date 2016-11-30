package jsonserv

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
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
	res := newResponse(mockWriter())
	m.Egress(nil, req, res)
	if c.egress != 1 {
		t.Fatalf("Egress incorrect: %d", c.egress)
	}
}

func TestLoggingMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewLoggingMiddleware(true)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse(mockWriter())
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
	res := newResponse(mockWriter())
	res.Err = errors.New("should be printed")
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar(StartTime).(time.Time).IsZero() {
		t.Fatal("Unexpected start time")
	}
	m.Egress(nil, req, res)
}

func TestMaxRequestSizeMiddleware_RequestVar_Is_Set(t *testing.T) {
	c := NewMaxRequestSizeMiddleware(5000)
	m := middlewares{c}
	req := newRequest(mockRequest())
	res := newResponse(mockWriter())
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
	res := newResponse(mockWriter())
	m.Ingress(nil, req, res)
	if req.GetMiddlewareVar("foo").(string) != "bar" {
		t.Fatal("Unexpected foo")
	}
	m.Egress(nil, req, res)
}

func TestGzipMiddleware_Wraps_Gzip_Accepted_ResponsesContentTypeIsSet(t *testing.T) {
	gz := NewGzipMiddleware()
	req := newRequest(mockRequest())
	res := newResponse(mockWriter())

	req.Header().Add(headerAcceptEncoding, headerAcceptEncodingGzip)
	gz.Ingress(nil, req, res)

	if _, ok := res.Writer.(*gzipWriter); !ok {
		t.Fatal("Writer not wrapped")
	}

}

func TestGzipMiddleware_Doesnt_wrap_non_Gzip_Accepted_ResponsesContentTypeIsSet(t *testing.T) {
	gz := NewGzipMiddleware()
	req := newRequest(mockRequest())
	res := newResponse(mockWriter())

	req.Header().Add(headerAcceptEncoding, "none")
	gz.Ingress(nil, req, res)

	if _, ok := res.Writer.(*gzipWriter); ok {
		t.Fatal("Writer wrapped")
	}

}

func TestGzipMiddleware_CompressesCorrectly(t *testing.T) {

	contents := []byte("hello, world!")

	compress := func(data []byte) []byte {
		buff := &bytes.Buffer{}
		compressor := gzip.NewWriter(buff)
		n, err := compressor.Write(data)
		if err != nil {
			t.Fatalf("Error compressing data: %v", err)
		}
		if n != len(data) {
			t.Fatal("Not all bytes compressed")
		}
		return buff.Bytes()
	}

	writer := mockWriter()
	gz := NewGzipMiddleware()
	req := newRequest(mockRequest())
	res := newResponse(writer)

	req.Header().Add(headerAcceptEncoding, headerAcceptEncodingGzip)
	gz.Ingress(nil, req, res)

	res.Writer.Write(contents)

	compressed := writer.Buffer.Bytes()
	expected := compress(contents)

	if reflect.DeepEqual(compressed, contents) {
		t.Fatal("Contents not compressed")
	}
	if !reflect.DeepEqual(compressed, expected) {
		t.Fatal("Unexpected gzipped content")
	}

}
