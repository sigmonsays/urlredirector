package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v9"
)

func (me *UrlHandler) CreateRedirect(w http.ResponseWriter, r *http.Request) error {
	rec := UrlRecord{}
	ctx := context.Background()

	if r.Method != "POST" {
		return me.sendError(w, r, InvalidRequest, "Invalid HTTP method: %s", r.Method)
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf, &rec)
	if err != nil {
		return me.sendError(w, r, Unknown, "Unmarshal request: %s", err)
	}
	now := time.Now()
	ts := now.Unix()
	rec.Ts = ts

	log.Printf("CreateRedirect %#v", rec)

	key := rec.Id

	if me.IsPathProtected(key) {
		return me.sendError(w, r, InvalidRequest, "ProtectedPath: %s", key)
	}

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
