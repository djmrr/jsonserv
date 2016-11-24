*This project is mostly just an exercise, but it works really well.*

# JSONServ
An easy way to spin up a JSON service

## Example

```go
package main

import (
	"log"
	"net/http"
	"errors"

	"github.com/explodes/jsonserv"
)

const (
        logRequestsAndResponse = true
        maxRequestSize = 1024 * 1024
        json500containsErrors = true
)


func main() {

    ctx := createAppContext() // app-level context such as database connections

	err := jsonserv.New().
		SetContext(ctx).
		AddMiddleware(jsonserv.NewLoggingMiddleware(logRequestsAndResponse)).
		AddMiddleware(jsonserv.NewMaxRequestSizeMiddleware(maxRequestSize)).
		AddMiddleware(jsonserv.NewDebugFlagMiddleware(json500containsErrors)).
		AddRoute(http.MethodGet, "Index", "/", contextWrap(indexView)).
		AddRoute(http.MethodGet, "Error", "/error", contextWrap(errorView)).
		Serve(":8888")
	if err != nil {
		log.Fatal(err)
	}
}

func contextWrap(f func(*AppContext, *jsonserv.Request, *jsonserv.Response)) jsonserv.View {
	return func(ctx interface{}, req *jsonserv.Request, res *jsonserv.Response) {
		f(ctx.(*AppContext), req, res)
	}
}

func indexView(ctx *AppContext, req *jsonserv.Request, res *jsonserv.Response) {
	res.Ok(map[string]interface{}{
		"hello":    "Hello, World!",
		"world":    true,
		"database": ctx.db,
		"producer": ctx.producer,
	})
}

func errorView(ctx *AppContext, req *jsonserv.Request, res *jsonserv.Response) {
	res.Error(errors.New("failed!!!"))
}
```