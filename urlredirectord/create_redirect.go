package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CreateResponse struct {
	Code    int
	Message string
}

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

	err = me.urlApi.CreateRedirect(ctx, &rec)
	if err != nil {
		return me.sendError(w, r, Unknown, "CreateRedirect: %s", key)
	}

	ret := &CreateResponse{}
	ret.Code = 0
	ret.Message = "created"
	buf2, _ := json.Marshal(ret)
	w.WriteHeader(200)
	w.Write(buf2)

	return nil
}
