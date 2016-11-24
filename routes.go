package jsonserv

import "fmt"

type route struct {
	name   string
	path   string
	method string
	view   View
}

func (r route) String() string {
	return fmt.Sprintf("%s=%s:%s", r.name, r.method, r.path)
}

type routes []*route

func (r *routes) Add(method, name, path string, view View) {
	*r = append(*r, &route{
		method: method,
		name:   name,
		path:   path,
		view:   view,
	})
}
