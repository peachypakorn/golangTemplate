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
	log.Debugln("StoreUpdate")

	var s Store
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	dBcontroller := Controller.RequestDBSession()
	if err := dBcontroller.StoreUpdate(bson.M{"_id": s.ID}, s); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func StoreFind(w http.ResponseWriter, r *http.Request) {
	//s := mux.Vars(r)
	//log.Errorln("mongo error:", s)
	key := r.FormValue("storename")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	selector := Store{StoreName:key}

	dBcontroller := Controller.RequestDBSession()
	result,err := dBcontroller.StoreFindByName(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	JSONResponse(w, http.StatusFound, result)


	w.WriteHeader(http.StatusOK)
}

func StoreDelete(w http.ResponseWriter, r *http.Request) {

}

func StoreGetAll(w http.ResponseWriter, r *http.Request) {
	dBcontroller := Controller.RequestDBSession()
	result,err := dBcontroller.StoreFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JSONResponse(w, http.StatusFound, result)
}

func StoreGetByID(w http.ResponseWriter, r *http.Request) {
	r.Header.Get("storeID")
}

