package jsonserv

import "testing"

func TestRoute_String(t *testing.T) {
	route := &route{}
	if route.String() == "" {
		t.Fatal("Bad string")
	}
}

func TestRoutes_Add(t *testing.T) {
	routes := make(routes, 0)
	routes.Add("GET", "Hello", "/", func(ctx interface{}, r *Request, out *Response) {})

	if len(routes) != 1 {
		t.Fatal("Route not added")
	}
}
