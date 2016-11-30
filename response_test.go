package jsonserv

import (
	"errors"
	"net/http"
	"testing"
)

func TestResponse_Done(t *testing.T) {
	res := newResponse(mockWriter())
	res.Done(5000, "abc")
	if res.Err != nil {
		t.Fatal("error set")
	}
	if res.Code != 5000 {
		t.Fatal("Code not set")
	}
	if res.Body != "abc" {
		t.Fatal("Body not set")
	}
}

func TestResponse_Empty(t *testing.T) {
	res := newResponse(mockWriter())
	res.Body = "abc"
	res.Empty(5000)
	if res.Err != nil {
		t.Fatal("error set")
	}
	if res.Code != 5000 {
		t.Fatal("Code not set")
	}
	if res.Body != nil {
		t.Fatal("Body not set")
	}
}

func TestResponse_Ok(t *testing.T) {
	res := newResponse(mockWriter())
	res.Ok("abc")
	if res.Err != nil {
		t.Fatal("error set")
	}
	if res.Code != http.StatusOK {
		t.Fatal("Code not set")
	}
	if res.Body != "abc" {
		t.Fatal("Body not set")
	}
}

func TestResponse_Err(t *testing.T) {
	res := newResponse(mockWriter())
	res.Error(errors.New("fail"))
	if res.Err == nil {
		t.Fatal("error not set")
	}
	if res.Code != http.StatusInternalServerError {
		t.Fatal("Code not set")
	}
	if res.Body != nil {
		t.Fatal("Body not set")
	}
}

func TestResponse_AddHeader(t *testing.T) {
	res := newResponse(mockWriter())
	res.AddHeader("content-type", "text/plain")
	if res.Writer.Header().Get("content-type") != "text/plain" {
		t.Fatal("Header not set")
	}
}
