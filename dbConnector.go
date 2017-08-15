package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Store struct {
	ID  	bson.ObjectId `json:"id" bson:"_id,omitempty"`
	StoreName 	string `json:"store_name" bson:"store_name,omitempty"`
	Branch    	string `json:"branch" bson:"branch,omitempty"`
	Phone     	string `json:"phone_num" bson:"phone_num,omitempty"`
	City		string `json:"city" bson:"city,omitempty"`
	Province	string `json:"province" bson:"province,omitempty"`
	//EnteringDay	Weekday `json:"entering_day" bson:"entering_day,omitempty"`
	//DayOff		Weekday `json:"day_off" bson:"day_off,omitempty"`
	//StartTime	timesHHMM `json:"sta" bson:"owner_id,omitempty"`
	//EndTime 	timesHHMM `json:"owner_id" bson:"owner_id,omitempty"`
	SellUnza	bool `json:"sell_unza" bson:"sell_unza,omitempty"`
	SellBio		bool `json:"sell_bio" bson:"sell_bio,omitempty"`
	DateAdd	time.Time `json:"date_add" bson:"date_add,omitempty"`
}

type DBcontroller struct {
	mongoSession *mgo.Session
}

var Controller DBcontroller

func (c DBcontroller) RequestDBSession() *DBcontroller {
	return &DBcontroller{mongoSession: getMongoSession() }
}

func (c DBcontroller) StoreInsert(store interface{}) error {

	session := c.mongoSession.Clone()
	defer session.Close()

	index := mgo.Index{
		Key:       []string {"store_name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	if err := session.DB("pcwutl").C("store").EnsureIndex(index); err != nil {
		return err
	}

	if err := session.DB("pcwutl").C("store").Insert(store); err != nil {
		return err
	}
	return nil
}

func (c DBcontroller) StoreFindByName(selector interface{}) (Store,error) {
	session := c.mongoSession.Clone()
	defer session.Close()

	var rtn Store
	if err := session.DB("pcwutl").C("store").Find(selector).One(&rtn); err != nil {
		return Store{},err
	}
	return rtn,nil
}


