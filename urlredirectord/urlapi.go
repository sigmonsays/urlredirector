package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
)

func NewUrlApi(rdb *redis.Client) *UrlApi {
	me := &UrlApi{
		rdb: rdb,
	}
	return me
}

type UrlApi struct {
	rdb *redis.Client
}

func (me *UrlApi) CreateRedirect(ctx context.Context, rec *UrlRecord) error {
	key := rec.Id
	_, err := me.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		pipeline.HSet(ctx, key, "id", rec.Id)
		pipeline.HSet(ctx, key, "url", rec.Url)
		pipeline.HSet(ctx, key, "ts", rec.Ts)
		//rdb.HSet(ctx, "key", "int", 123)
		//rdb.HSet(ctx, "key", "bool", 1)
		return nil
	})
	return err
}

func (me *UrlApi) GetRedirect(ctx context.Context, key string) (*UrlRecord, error) {
	var rec UrlRecord
	err := me.rdb.HGetAll(ctx, key).Scan(&rec)
	if err != nil {
		log.Printf("ERROR: HGetAll %s: %s", key, err)
		return nil, err
	}
	return &rec, nil
}
