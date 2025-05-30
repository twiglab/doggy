package hx

import (
	"encoding/json"
	"net/http"
)

func JsonTo(code int, resp any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}

func BindAndClose(r *http.Request, p any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(p); err != nil {
		return err
	}
	return nil
}

func Bind(r *http.Request, p any) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(p); err != nil {
		return err
	}
	return nil
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
