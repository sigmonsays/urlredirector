package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/sethvargo/go-envconfig"
	"github.com/sigmonsays/urlredirector/static"
)

type Config struct {
	HttpPort  int    `env:"HTTP_PORT,default=8080"`
	RedisHost string `env:"REDIS_HOST,default=localhost"`
	RedisPort int    `env:"REDIS_PORT,default=6379"`

	Auth struct {
		Username string `env:"AUTH_USERNAME"`
		Password string `env:"AUTH_PASSWORD"`
	} `env:"auth"`
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
func run() error {
	cfg := Config{}
	ctx := context.Background()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	log.Printf("app config %+v", cfg)

	redisAddr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	log.Printf("connecting to redis at %s", redisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	urlApi := NewUrlApi(rdb)

	rec := &UrlRecord{}
	rec.Id = "/welcome/index.html"
	rec.Url = "/welcome/index.html"
	urlApi.CreateRedirect(ctx, rec)

	handler := &UrlHandler{}
	handler.urlApi = urlApi
	handler.protectedPaths = []string{
		"/api",
	}

	eh := &ErrorWrapper{}

	mx := http.NewServeMux()

	haveAuth := cfg.Auth.Username != "" && cfg.Auth.Password != ""
	var createHandler http.HandlerFunc
	createHandler = eh.WrapHandler(handler.CreateRedirect)
	if haveAuth {
		basicAuth := &BasicAuth{}
		basicAuth.Username = cfg.Auth.Username
		basicAuth.Password = cfg.Auth.Password
		createHandler = basicAuth.AuthWrapper(createHandler)
		log.Printf("Setting up BasicAuth on server")
	}
	mx.HandleFunc("/api/create", createHandler)

	mx.Handle("/welcome/", http.StripPrefix("/welcome/", http.FileServer(http.FS(static.Files))))
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
