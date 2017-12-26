package main

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
)
var mongoSession *mgo.Session

func getMongoSession() *mgo.Session {
	return mongoSession

}

func startMongo() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Errorln("error mongo:", err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetPoolLimit(3000)
	mongoSession = session
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.Kitchen})
	log.SetLevel(log.DebugLevel)
	startMongo()
	n := negroni.New()
	n.Use(negroni.NewRecovery())

	n.UseHandler(NewRouter())
	n.Run(":4444")
}
