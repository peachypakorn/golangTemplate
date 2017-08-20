package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type ResponseNormal struct {
	Code          string  `json:"responseCode"`
	Description   string  `json:"responseDescription"`
	StoreID		  string  `json:"storeID"`
}

func JSONResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Errorln("JSON encoder error:", err)
	}
}
