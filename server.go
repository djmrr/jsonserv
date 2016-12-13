package jsonserv

import (
	"net/http"

	"github.com/gorilla/mux"
	"net"
	"time"
	"errors"
)

const (
	contentTypeHeader = "Content-type"
	contentTypeJson = "application/json"
	emptyBody = "{}"
)

type JsonServer struct {
	App         interface{}
	routes      routes
	Middlewares middlewares
	Listener    *net.TCPListener
}

func New() *JsonServer {
	return &JsonServer{
		routes:      make(routes, 0, 16),
		Middlewares: make(middlewares, 0, 2),
	}
}

func (s *JsonServer) AddRoute(method, name, path string, view View) *JsonServer {
	s.routes.Add(method, name, path, view)
	return s
}

func (s *JsonServer) AddMiddleware(middleware Middleware) *JsonServer {
	s.Middlewares = append(s.Middlewares, middleware)
	return s
}

func (s *JsonServer) SetApp(app interface{}) *JsonServer {
	s.App = app
	return s
}

func (s *JsonServer) Serve(addr string) error {
	router := s.createRouter()
	server := &http.Server{Addr: addr, Handler: router}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.Listener = ln.(*net.TCPListener)
	return server.Serve(tcpKeepAliveListener{s.Listener})
}

func (s *JsonServer) Close() error {
	if s.Listener == nil {
		return errors.New("Server not listening")
	}
	old := s.Listener
	s.Listener = nil
	return old.Close()
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
	return s.newHandler("NotFound", func(app interface{}, r *Request, out *Response) {
		out.Empty(http.StatusNotFound)
	})
}

func (s *JsonServer) newHandler(name string, view View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		res := newWrappedResponse(w)
		defer func() {
			res.Writer.Close()
		}()

		s.Middlewares.Ingress(s.App, req, res)
		view(s.App, req, res)
		s.Middlewares.Egress(s.App, req, res)

		respond(req, res)
	})
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}