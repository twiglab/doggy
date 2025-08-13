package hx

import (
	"encoding/json/v2"
	"net/http"
)

func JsonTo(code int, resp any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.MarshalWrite(w, resp)
}

// Deprecated: use Bind
func BindAndClose(r *http.Request, p any) error {
	defer r.Body.Close()
	return json.UnmarshalRead(r.Body, p)
}

func Bind(r *http.Request, p any) error {
	return json.UnmarshalRead(r.Body, p)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
