package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v9"
)

type UrlRecord struct {
	Id  string `redis:"id"`
	Url string `redis:"url"`
	Ts  int64  `redis:"ts"`
}

type UrlHandler struct {
	rdb            *redis.Client
	protectedPaths []string
}

type UrlResponse struct {
	Error string `json:"error,omitempty"`
	Url   string `json:"url,omitempty"`
}

func (me *UrlHandler) sendError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) error {
	msg := fmt.Sprintf(s, args...)
	log.Printf("sendError %s: %s", r.URL, msg)

	res := &UrlResponse{}
	res.Error = msg
	buf, _ := json.MarshalIndent(res, "", " ")

	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(buf)
	return nil
}

func (me *UrlHandler) IsPathProtected(path string) bool {
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
