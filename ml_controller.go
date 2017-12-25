package main

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type RecommendedProduct struct{
	Data []Product `json:"data" bson:"data,omitempty"`
}

type Product struct {
	Pid string `json:"pid" bson:"pid,omitempty"`
	Score string `json:"score" bson:"score,omitempty"`

}

func getRecommendedProducts(w http.ResponseWriter, r *http.Request) {
	//s := mux.Vars(r)
	//log.Errorln("mongo error:", s)
	key := mux.Vars(r)
	if key == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()


	selector := Cus{Uid:key["uid"]}

	dBcontroller := Controller.RequestDBSession()
	result,err := dBcontroller.GetRecommendedProducts(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debug(result)
	JSONResponse(w, http.StatusFound, "eiei")


	w.WriteHeader(http.StatusOK)
}
