package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func CheckDestAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		// BeaconID      string  `json:"beaconId"`
		DestBankName string `json:"destBankName,omitempty"`
		DestAccNo    string `json:"destAccNo,omitempty"`
		DestPhoneNo  string `json:"destPhoneNo,omitempty"`
	}

	type res struct {
		Code              string `json:"responseCode"`
		Description       string `json:"responseDescription"`
		DestAccountName   string `json:"destAccName"`
		DestBankName      string `json:"destBankName"`
		DestAccountNumber string `json:"destAccNo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	var selector bson.M
	if req.DestBankName != "" && req.DestAccNo != "" {
		selector = bson.M{
			"bank_name":      req.DestBankName,
			"account_number": req.DestAccNo,
		}
	} else {
		selector = bson.M{
			"phone_number": req.DestPhoneNo,
		}
	}

	asset, err := db.AssetFind(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	user, err := db.CustomerFind(bson.M{"_id": asset.OwnerID})
	var name string
	if err != nil {
		log.Errorln("mongo error:", err)
		name = ""
	} else {
		name = user.FirstName + " " + user.LastName
	}

	log.Println(asset)
	JSONResponse(w, http.StatusFound, res{
		Code:              "00",
		Description:       "success",
		DestBankName:      asset.BankName,
		DestAccountNumber: asset.AccNumber,
		DestAccountName:   name,
	})
}
