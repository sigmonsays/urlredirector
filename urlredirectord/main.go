package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HttpPort  int    `env:"HTTP_PORT"`
	RedisHost string `env:"REDIS_HOST"`
	RedisPort int    `env:"REDIS_PORT"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
func run() error {
	cfg := &Config{}
	cfg.HttpPort = 8080
	ctx := context.Background()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	handler := &UrlHandler{}
	handler.rdb = rdb

	eh := &ErrorWrapper{}

	mx := http.NewServeMux()
	mx.HandleFunc("/api/create", eh.WrapHandler(handler.CreateRedirect))
	mx.HandleFunc("/", eh.WrapHandler(handler.GetRedirect))

	srv := &http.Server{}
	srv.Addr = fmt.Sprintf(":%d", cfg.HttpPort)
	srv.Handler = mx

	log.Printf("Listening at port %d", cfg.HttpPort)
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

type UrlRecord struct {
	Id  string `redis:"id"`
	Url string `redis:"url"`
	Ts  int64  `redis:"ts"`
}

type UrlHandler struct {
	rdb *redis.Client
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

func (me *UrlHandler) CreateRedirect(w http.ResponseWriter, r *http.Request) error {
	rec := UrlRecord{}
	ctx := context.Background()

	if r.Method != "POST" {
		return me.sendError(w, r, "Invalid HTTP method: %s", r.Method)
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf, &rec)
	if err != nil {
		return err
	}
	now := time.Now()
	ts := now.Unix()
	rec.Ts = ts

	log.Printf("CreateRedirect %#v", rec)

	key := rec.Id

	_, err = me.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, key, "id", rec.Id)
		rdb.HSet(ctx, key, "url", rec.Url)
		rdb.HSet(ctx, key, "ts", rec.Ts)
		//rdb.HSet(ctx, "key", "int", 123)
		//rdb.HSet(ctx, "key", "bool", 1)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *UrlHandler) GetRedirect(w http.ResponseWriter, r *http.Request) error {
	rec := UrlRecord{}
	ctx := context.Background()

	key := r.URL.Path
	log.Printf("GetRedirect %s", key)

	err := me.rdb.HGetAll(ctx, key).Scan(&rec)
	if err != nil {
		log.Printf("ERROR: HGetAll %s: %s", key, err)
		return nil
	}

	log.Printf("GetRedirect got %#v", rec)
	if rec.Id == "" {
		return me.sendError(w, r, "No such redirect: %s", key)
	}

	w.Header().Set("Location", rec.Url)
	w.WriteHeader(302)

	return nil
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
