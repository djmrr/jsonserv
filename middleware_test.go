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

func compress(b []byte) ([]byte, error) {
	buff := &bytes.Buffer{}
	gz := gzip.NewWriter(buff)
	n, err := gz.Write(b)
	if err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}
	if n != len(b) {
		return nil, errors.New("Not all bytes compressed")
	}
	return buff.Bytes(), nil
}

func decompress(b []byte) ([]byte, error) {
	buff := bytes.NewReader(b)
	gz, err := gzip.NewReader(buff)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return ioutil.ReadAll(gz)
}

func TestGzipWriter_CompressesCorrectly(t *testing.T) {

	contents := []byte("hello, world!")
	expected, err := compress(contents)
	if err != nil {
		t.Fatalf("Error compressing contents: %v", err)
	}

	writer := mockWriter()
	gzWriter := &gzipWriter{
		writer: writer,
		gz:     gzip.NewWriter(writer),
	}

	n, err := gzWriter.Write(contents)
	gzWriter.Close()
	if err != nil {
		t.Fatalf("Unable to compress contents: %v", err)
	}
	if n != len(contents) {
		t.Fatalf("Unable to write all contents: %d/%d", n, len(contents))
	}

	results := writer.Buffer.Bytes()
	if reflect.DeepEqual(results, contents) {
		t.Fatal("Contents not compressed")
	}
	if !reflect.DeepEqual(results, expected) {
		t.Fatal("Unexpected gzipped content")
	}

	decompressed, err := decompress(results)
	if err != nil {
		t.Fatalf("Error decompressing results: %v", err)
	}
	if !reflect.DeepEqual(contents, decompressed) {
		t.Fatal("Contents did not decompress to original")
	}

}

func TestGzipMiddleware_RoundTrip(t *testing.T) {
	contents := []byte("hello, world!")
	expected, err := compress(contents)
	if err != nil {
		t.Fatalf("Error compressing contents: %v", err)
	}

	writer := mockWriter()
	gz := NewGzipMiddleware()
	req := newRequest(mockRequest())
	res := newResponse(writer)
	req.Header().Add(headerAcceptEncoding, headerAcceptEncodingGzip)
	gz.Ingress(nil, req, res)
	gz.Egress(nil, req, res)

	res.Writer.Write(contents)
	// Close usually happens during response but we need to force it here to check output
	res.Writer.(*gzipWriter).Close()

	results := writer.Buffer.Bytes()
	if !reflect.DeepEqual(results, expected) {
		t.Fatal("Unexpected results")
	}

	decompressed, err := decompress(results)
	if err != nil {
		t.Fatalf("Error decompressing results: %v", err)
	}
	if !reflect.DeepEqual(contents, decompressed) {
		t.Fatal("Contents did not decompress to original")
	}

}
