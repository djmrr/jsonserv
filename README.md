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
	"database/sql"

	"github.com/explodes/jsonserv"
	"github.com/explodes/ezconfig/producer"
)

const (
        logRequestsAndResponse = true
        maxRequestSize = 1024 * 1024
        json500containsErrors = true
)

// App is service-wide state such as database connections,
// configuration variables, info about remote services, etc.
type App struct{ /**/ 
        db *sql.DB
        producer producer.Producer
}
func createApp() *App{ /**/ return &App{} }

func main() {

    app := createApp() // app-level context such as database connections

	err := jsonserv.New().
		SetApp(app).
		AddMiddleware(jsonserv.NewLoggingMiddleware(logRequestsAndResponse)).
		AddMiddleware(jsonserv.NewMaxRequestSizeMiddleware(maxRequestSize)).
		AddMiddleware(jsonserv.NewDebugFlagMiddleware(json500containsErrors)).
		AddRoute(http.MethodGet, "Index", "/", appWrap(indexView)).
		AddRoute(http.MethodGet, "Error", "/error", appWrap(errorView)).
		Serve(":8888")
	if err != nil {
		log.Fatal(err)
	}
}

// appWrap wraps a view to pre-type-assert the app interface{} parameter into our App's type
func appWrap(f func(*App, *jsonserv.Request, *jsonserv.Response)) jsonserv.View {
	return func(app interface{}, req *jsonserv.Request, res *jsonserv.Response) {
		f(app.(*App), req, res)
	}
}

func indexView(app *App, req *jsonserv.Request, res *jsonserv.Response) {
	res.Ok(map[string]interface{}{
		"hello":    "Hello, World!",
		"world":    true,
		"database": app.db,
		"producer": app.producer,
	})
}

func errorView(app *App, req *jsonserv.Request, res *jsonserv.Response) {
	res.Error(errors.New("failed!!!"))
}
```