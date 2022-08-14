package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type UrlRecord struct {
	Id  string `redis:"id"`
	Url string `redis:"url"`
	Ts  int64  `redis:"ts"`
}

type UrlHandler struct {
	protectedPaths []string
	urlApi         *UrlApi
}

type UrlResponse struct {
	Reason string `json:"reason,omitempty"`
	Error  string `json:"error,omitempty"`
	Url    string `json:"url,omitempty"`
}

func (me *UrlHandler) sendError(w http.ResponseWriter, r *http.Request, ec ErrorClass, s string, args ...interface{}) error {

	res := &UrlResponse{}

	if ec == Unknown {
		msg := fmt.Sprintf(s, args...)
		res.Error = msg
	} else {
		ae := ec.Sprintf(s, args...)
		res.Reason = ae.Class.String()
		res.Error = ae.Message
	}
	log.Printf("sendError %s: %+v", r.URL, res)

	buf, _ := json.MarshalIndent(res, "", " ")

	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(buf)
	return nil
}

func (me *UrlHandler) IsPathProtected(path string) bool {
	if path == "" || path == "/" {
		return true
	}
	for _, p := range me.protectedPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

type ErrorWrapper struct {
}

func (me *ErrorWrapper) WrapHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Printf("ErrorWrapper %s: %s", r.URL, err)
		}
	}
}
