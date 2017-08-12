package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func ListBankName(w http.ResponseWriter, r *http.Request) {
	log.Debugln("ListBankName")

	type DestBankName struct {
		Name string `json:"destBankName"`
	}

	type res struct {
		Code        string         `json:"responseCode"`
		Description string         `json:"responseDescription"`
		Bank        []DestBankName `json:"bank"`
	}

	db := RP.Request(r)
	a, err := db.AssetFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	unique := map[string]int{}

	for _, asset := range a {
		unique[asset.BankName] = 1
	}

	bankNames := []DestBankName{}
	for name := range unique {
		var entry DestBankName
		entry.Name = name
		bankNames = append(bankNames, entry)
	}

	JSONResponse(w, http.StatusFound, res{
		Code:        "00",
		Description: "success",
		Bank:        bankNames,
	})
}
