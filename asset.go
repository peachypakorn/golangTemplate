package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func AssetCreate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AssetCreate")

	var req Asset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.AssetInsert(req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AssetUpdate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AssetUpdate")

	var req Asset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.AssetUpdate(bson.M{"_id": req.ID}, req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AssetGetAll(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AssetGetAll")

	db := RP.Request(r)
	a, err := db.AssetFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JSONResponse(w, http.StatusFound, a)
}

func AliasCreate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AliasCreate")

	var req Alias
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Debugln(req)

	db := RP.Request(r)
	if err := db.AliasInsert(req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AliasUpdate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AliasUpdate")

	var req Alias
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.AliasUpdate(bson.M{"_id": req.ID}, req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AliasGetAll(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AliasGetAll")

	db := RP.Request(r)
	a, err := db.AliasFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JSONResponse(w, http.StatusFound, a)
}
