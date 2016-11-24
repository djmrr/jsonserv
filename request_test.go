package jsonserv

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req == nil {
		t.Fatal("nil request")
	}
}

func TestRequest_Method(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.Method() != http.MethodGet {
		t.Fatal("Unexpected method")
	}
}

func TestRequest_URL(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.URL() != exampleUrl {
		t.Fatal("Unexepected URL")
	}
}

func TestRequest_GetMiddlewareVar(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	req.SetMiddlewareVar("foo", "bar")
	if req.GetMiddlewareVar("foo").(string) != "bar" {
		t.Fatal("Bad foo")
	}
}

func TestRequest_GetOptionalMiddlewareVar(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.GetOptionalMiddlewareVar("foo", "bar").(string) != "bar" {
		t.Fatal("Bad default")
	}
}

func TestRequest_String(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.String() == "" {
		t.Fatal("Bad string")
	}
}

func TestRequest_GetPathVars(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.GetPathVars() != nil {
		t.Fatal("Vars should be nil without gorilla mux")
	}
}

func TestRequest_GetPathVar(t *testing.T) {
	r := mockRequest()
	req := newRequest(r)
	if req.GetPathVar("foo", "bar") != "bar" {
		t.Fatal("Unexpected var")
	}
}

func TestRequest_ParseBody_no_limit(t *testing.T) {
	type requestBody struct {
		Foo string `json:"foo"`
	}
	body := &requestBody{}

	r := mockRequest()
	req := newRequest(r)

	if err := req.ParseBody(body); err != nil {
		t.Fatal(err)
	}
	if body.Foo != "bar" {
		t.Fatal("Unexpected json")
	}
}

func TestRequest_ParseBody_limited(t *testing.T) {
	type requestBody struct {
		Foo string `json:"foo"`
	}
	body := &requestBody{}

	r := mockRequest()
	req := newRequest(r)
	req.SetMiddlewareVar(MaxBodySize, int64(5000))

	if err := req.ParseBody(body); err != nil {
		t.Fatal(err)
	}
	if body.Foo != "bar" {
		t.Fatal("Unexpected json")
	}
}

func TestRequest_ParseBody_too_small(t *testing.T) {
	type requestBody struct {
		Foo string `json:"foo"`
	}
	body := &requestBody{}

	r := mockRequest()
	req := newRequest(r)
	req.SetMiddlewareVar(MaxBodySize, int64(3))

	err := req.ParseBody(body)
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "Request body too large" {
		t.Fatalf("Unexpected error: %v", err)
	}
}
