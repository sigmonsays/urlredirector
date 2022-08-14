package main

import (
	"context"
	"log"
	"net/http"
)

func (me *UrlHandler) GetRedirect(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	key := r.URL.Path
	log.Printf("GetRedirect %s", key)

	rec, err := me.urlApi.GetRedirect(ctx, key)
	if err != nil {
		log.Printf("ERROR: HGetAll %s: %s", key, err)
		return nil
	}

	log.Printf("GetRedirect got %#v", rec)
	if rec.Id == "" {
		return me.sendError(w, r, NoSuchRedirect, "No such redirect: %s", key)
	}

	w.Header().Set("Location", rec.Url)
	w.WriteHeader(302)

	return nil
}
