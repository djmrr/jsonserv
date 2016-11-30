package jsonserv

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func respond(req *Request, res *Response) {
	var err error
	if res.Err != nil {
		err = writeError(req, res)
	} else {
		err = writeBody(res)
	}
	if err != nil {
		log.Printf("Error rendering %s: %v", req.URL(), err)
	}
}

func writeError(req *Request, res *Response) error {
	body := make(map[string]interface{})
	if req.GetOptionalMiddlewareVar(DebugFlag, false).(bool) {
		body["error"] = res.Err.Error()
	}
	return write(res.Writer, http.StatusInternalServerError, body)
}

func writeBody(res *Response) error {
	return write(res.Writer, res.Code, res.Body)
}

func write(w http.ResponseWriter, code int, body interface{}) error {
	w.Header().Add(contentTypeHeader, contentTypeJson)
	w.WriteHeader(code)
	if body == nil {
		return writeEmptyBody(w)
	} else {
		enc := json.NewEncoder(w)
		return enc.Encode(body)
	}
}

func writeEmptyBody(w http.ResponseWriter) error {
	if n, err := w.Write([]byte(emptyBody)); err != nil || n != len(emptyBody) {
		if err != nil {
			return err
		} else {
			return errors.New("Empty body not fully written")
		}
	} else {
		return nil
	}
}
