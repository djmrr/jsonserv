package jsonserv

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, req *Request, res *Response) {
	var err error
	if res.Err != nil {
		err = writeError(w, req, res)
	} else {
		err = writeBody(w, res)
	}
	if err != nil {
		log.Printf("Error rendering %s: %v", req.URL(), err)
	}
}

func writeError(w http.ResponseWriter, req *Request, res *Response) error {
	body := make(map[string]interface{})
	if req.GetOptionalMiddlewareVar(DebugFlag, false).(bool) {
		body["error"] = res.Err.Error()
	}
	return write(w, http.StatusInternalServerError, res.headers, body)
}

func writeBody(w http.ResponseWriter, res *Response) error {
	return write(w, res.Code, res.headers, res.Body)
}

func write(w http.ResponseWriter, code int, headers map[string]string, body interface{}) error {
	w.Header().Add(contentTypeHeader, contentTypeJson)
	w.WriteHeader(code)
	if headers != nil {
		for k, v := range headers {
			w.Header().Add(k, v)
		}
	}
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
