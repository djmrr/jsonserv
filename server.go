package jsonserv

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	contentTypeHeader = "Content-type"
	contentTypeJson   = "application/json"
	emptyBody         = "{}"
)

type JsonServer struct {
	Context     interface{}
	routes      routes
	middlewares middlewares
}

func New() *JsonServer {
	return &JsonServer{
		routes:      make(routes, 0, 16),
		middlewares: make(middlewares, 0, 2),
	}
}

func (s *JsonServer) AddRoute(method, name, path string, view View) *JsonServer {
	s.routes.Add(method, name, path, view)
	return s
}

func (s *JsonServer) AddMiddleware(middleware Middleware) *JsonServer {
	s.middlewares = append(s.middlewares, middleware)
	return s
}

func (s *JsonServer) SetContext(ctx interface{}) *JsonServer {
	s.Context = ctx
	return s
}

func (s *JsonServer) Serve(bind string) error {
	router := s.createRouter()
	return http.ListenAndServe(bind, router)
}

func (s *JsonServer) createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range s.routes {
		handler := s.newHandler(route.name, route.view)
		router.Methods(route.method).
			Path(route.path).
			Name(route.name).
			Handler(handler)
	}
	router.NotFoundHandler = s.newNotFoundHandler()
	return router
}

func (s *JsonServer) newNotFoundHandler() http.Handler {
	return s.newHandler("NotFound", func(ctx interface{}, r *Request, out *Response) {
		out.Empty(http.StatusNotFound)
	})
}

func (s *JsonServer) newHandler(name string, view View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newResponse()

		s.middlewares.Ingress(s.Context, req, res)
		view(s.Context, req, res)
		s.middlewares.Egress(s.Context, req, res)

		respond(w, req, res)
	})
}
