package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func CustomerGetAll(w http.ResponseWriter, r *http.Request) {
	log.Debugln("CustomerGetAll")

	db := RP.Request(r)
	c, err := db.CustomerFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JSONResponse(w, http.StatusFound, c)
}

func CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("CustomerUpdate")

	var req Customer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.CustomerUpdate(bson.M{"_id": req.ID}, req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CustomerCreate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("CustomerCreate")

	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.CustomerInsert(c); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
