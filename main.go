package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
)

func SetDB(r *http.Request, key string, val *mgo.Database) {
	context.Set(r, key, val)
}

func MongoMiddleware() negroni.HandlerFunc {
	bs, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Errorln("error mongo:", err)
	}
	bs.SetMode(mgo.Monotonic, true)

	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		beacon := bs.Clone()
		defer beacon.Close()
		SetDB(r, "beacon", beacon.DB("beacon"))
		next(w, r)
	})
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.Kitchen})
	log.SetLevel(log.DebugLevel)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(MongoMiddleware())

	//mux := http.NewServeMux()
	//mux.HandleFunc("/customers", CustomerGetAll)
	//mux.HandleFunc("/customer/update", CustomerUpdate)
	//mux.HandleFunc("/customer/create", CustomerCreate)
	//mux.HandleFunc("/customer/asset/create", CustomerInsertAsset)
	//mux.HandleFunc("/transfer", TransferBalance)
	//mux.HandleFunc("/asset/update", AssetUpdateBankAccount)
	//mux.HandleFunc("/asset/create", AssetInsertBankAccount)

	//n.UseHandler(mux)

	n.UseHandler(NewRouter())
	n.Run(":4444")
}
