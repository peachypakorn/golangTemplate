package main

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"math/rand"
	"strconv"
)

type RecommendedProduct struct {
	Data []Product `json:"data" bson:"data,omitempty"`
}

type Product struct {
	Pid   string `json:"pid" bson:"pid,omitempty"`
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
	currentUid := key["uid"]
	if (len(currentUid) < 1) {
		currentUid = strconv.Itoa(rand.Intn(10))
	} else {
		currentUid = currentUid[len(currentUid) - 2:]
	}
	//log.Debug(currentUid)
	selector := Cus{Uid:key["uid"]}

	dBcontroller := Controller.RequestDBSession()
	results, err := dBcontroller.GetRecommendedProducts(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if len(results) == 0 {
		log.Errorln("not Found")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//log.Debug(result)

	var response []Product
	for _ , result := range results{
		response = append(response,Product{
			Pid:result.Pid,
			Score:result.Score,
		})
	}

	JSONResponse(w, http.StatusFound, RecommendedProduct{Data:response})

	w.WriteHeader(http.StatusOK)
}

func addMockProducts(w http.ResponseWriter, r *http.Request) {
	var data []Cus
	for i := 0; i < 100; i++ {
		// Display integer.
		products := rand.Intn(99)
		for j := 0; j < products; j++ {
			cus := Cus{Uid:"" + strconv.Itoa(i)}
			cus.Pid = strconv.Itoa(rand.Intn(999999))
			cus.Score = strconv.Itoa(rand.Intn(100))
			data = append(data, cus)
		}
	}

	dBcontroller := Controller.RequestDBSession()
	err := dBcontroller.AddProduct(data)
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	JSONResponse(w, http.StatusFound, "add mock data success")

	w.WriteHeader(http.StatusOK)
}
