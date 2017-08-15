package main

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"time"
	"gopkg.in/mgo.v2/bson"
)



func StoreCreate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("StoreCreate")

	var s Store
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	s.DateAdd = time.Now()

	dBcontroller := Controller.RequestDBSession()
	s.ID = bson.NewObjectId()

	if err := dBcontroller.StoreInsert(&s); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	JSONResponse(w, http.StatusCreated, ResponseNormal{
		"00",
		"create store success",
		s.ID.Hex(),
	})

}

func StoreUpdate(w http.ResponseWriter, r *http.Request) {

}

func StoreDelete(w http.ResponseWriter, r *http.Request) {

}

func StoreGetAll(w http.ResponseWriter, r *http.Request) {

}

func StoreGetByID(w http.ResponseWriter, r *http.Request) {
	r.Header.Get("storeID")
}