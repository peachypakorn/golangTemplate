package main

import (
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Asset struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	OwnerID     bson.ObjectId `json:"owner_id" bson:"owner_id,omitempty"`
	Type        string        `json:"type" bson:"type,omitempty"`
	BankName    string        `json:"bank_name" bson:"bank_name,omitempty"`
	AccNumber   string        `json:"account_number" bson:"account_number,omitempty"`
	Currency    string        `json:"currency" bson:"currency,omitempty"`
	Balance     float64       `json:"balance" bson:"balance,omitempty"`
	Status      string        `json:"status" bson:"status,omitempty"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number,omitempty"`
}

type Customer struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
}

type Alias struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id,omitempty"`
	AssetID bson.ObjectId `json:"asset_id" bson:"asset_id,omitempty"`
	Alias   string        `json:"alias" bson:"alias,omitempty"`
}

type Bill struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	BillType        string        `json:"billType,omitempty" bson:"billType,omitempty"`
	BillService     string        `json:"billService,omitempty" bson:"billService,omitempty"`
	ServiceProvider string        `json:"serviceProvider,omitempty" bson:"serviceProvider,omitempty"`
	Ref1            string        `json:"REF1,omitempty" bson:"REF1,omitempty"`
	Ref2            string        `json:"REF2,omitempty" bson:"REF2,omitempty"`
	BillAccNo       string        `json:"billAccNo,omitempty" bson:"billAccNo,omitempty"`
	BillRef1        string        `json:"billREF1,omitempty" bson:"billREF1,omitempty"`
	BillRef2        string        `json:"billREF2,omitempty" bson:"billREF2,omitempty"`
	DueDate         time.Time     `json:"dueDate,omitempty" bson:"dueDate,omitempty"`
	BilledAmount    float64       `json:"billedAmount,omitempty" bson:"billedAmount,omitempty"`
	PaidAmount      float64       `json:"paidAmount,omitempty" bson:"paidAmount,omitempty"`
	MinAmount       float64       `json:"minAmount,omitempty" bson:"minAmount,omitempty"`
	Currency        string        `json:"currency,omitempty" bson:"currency,omitempty"`
}

type Repo struct {
	bs *mgo.Database
}

var RP Repo

func GetDB(r *http.Request, name string) *mgo.Database {
	if rv := context.Get(r, name); rv != nil {
		return rv.(*mgo.Database)
	}
	return nil
}

func (r Repo) Request(req *http.Request) *Repo {
	return &Repo{bs: GetDB(req, "beacon")}
}

func (r Repo) CustomerFindAll() ([]Customer, error) {
	var rtn []Customer
	if err := r.bs.C("customers").Find(nil).All(&rtn); err != nil {
		return []Customer{}, err
	}
	return rtn, nil
}

func (r Repo) CustomerFind(selector interface{}) (Customer, error) {
	var rtn Customer
	if err := r.bs.C("customers").Find(selector).One(&rtn); err != nil {
		return Customer{}, err
	}
	return rtn, nil
}

func (r Repo) CustomerUpdate(selector interface{}, update interface{}) error {
	if err := r.bs.C("customers").Update(selector, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r Repo) CustomerInsert(customer interface{}) error {
	if err := r.bs.C("customers").Insert(customer); err != nil {
		return err
	}
	return nil
}

func (r Repo) AssetFindAll() ([]Asset, error) {
	var rtn []Asset
	if err := r.bs.C("assets").Find(nil).All(&rtn); err != nil {
		return []Asset{}, err
	}
	return rtn, nil
}

func (r Repo) AssetFind(selector interface{}) (Asset, error) {
	var rtn Asset
	if err := r.bs.C("assets").Find(selector).One(&rtn); err != nil {
		return Asset{}, err
	}
	return rtn, nil
}

func (r Repo) AssetFindByID(assetID string) (Asset, error) {
	var rtn Asset
	if err := r.bs.C("assets").FindId(bson.ObjectIdHex(assetID)).One(&rtn); err != nil {
		return Asset{}, err
	}
	return rtn, nil
}

func (r Repo) AssetUpdate(selector interface{}, update interface{}) error {
	if err := r.bs.C("assets").Update(selector, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r Repo) AssetInsert(asset interface{}) error {
	if err := r.bs.C("assets").Insert(asset); err != nil {
		return err
	}
	return nil
}

// alias
func (r Repo) AliasFindAll() ([]Alias, error) {
	var rtn []Alias
	if err := r.bs.C("aliases").Find(nil).All(&rtn); err != nil {
		return []Alias{}, err
	}
	return rtn, nil
}

func (r Repo) AliasFind(selector interface{}) (Alias, error) {
	var rtn Alias
	if err := r.bs.C("aliases").Find(selector).One(&rtn); err != nil {
		return Alias{}, err
	}
	return rtn, nil
}

func (r Repo) AliasUpdate(selector interface{}, update interface{}) error {
	if err := r.bs.C("aliases").Update(selector, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r Repo) AliasInsert(alias interface{}) error {
	if err := r.bs.C("aliases").Insert(alias); err != nil {
		return err
	}
	return nil
}

func (r Repo) TransactionInsert(t interface{}) error {
	if err := r.bs.C("transactions").Insert(t); err != nil {
		return err
	}
	return nil
}

func (r Repo) TransactionFind(selector interface{}) (Transaction, error) {
	var rtn Transaction
	if err := r.bs.C("transactions").Find(selector).One(&rtn); err != nil {
		return Transaction{}, err
	}
	return rtn, nil
}

func (r Repo) TransactionFindByType(ownerID, transactionType string, limit int) ([]Transaction, error) {
	var rtn []Transaction
	s := bson.M{"owner_id": bson.ObjectIdHex(ownerID), "type": transactionType}
	err := r.bs.C("transactions").Find(s).Sort("-_id").Limit(limit).All(&rtn)
	if err != nil {
		return []Transaction{}, err
	}
	return rtn, nil
}

func (r Repo) TransactionFindFavorites(ownerID, transactionType string) ([]Transaction, error) {
	var rtn []Transaction
	s := bson.M{"favorite": bson.M{"$exists": true}, "owner_id": bson.ObjectIdHex(ownerID), "type": transactionType}
	err := r.bs.C("transactions").Find(s).Sort("-_id").All(&rtn)
	if err != nil {
		return []Transaction{}, err
	}
	return rtn, nil
}

func (r Repo) TransactionUpdate(selector interface{}, update interface{}) error {
	if err := r.bs.C("transactions").Update(selector, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r Repo) TransactionAddFavorite(transactionID, alias string) error {
	update := bson.M{"favorite": alias}
	if err := r.TransactionUpdate(bson.M{"_id": bson.ObjectIdHex(transactionID)}, update); err != nil {
		return err
	}
	return nil
}

func (r Repo) TransactionUnsetField(transactionID, field string) error {
	update := bson.M{"$unset": bson.M{field: ""}}
	if err := r.bs.C("transactions").UpdateId(bson.ObjectIdHex(transactionID), update); err != nil {
		return err
	}
	return nil
}

func (r Repo) BillFind(selector interface{}) (Bill, error) {
	var rtn Bill
	if err := r.bs.C("bills").Find(selector).One(&rtn); err != nil {
		return Bill{}, err
	}
	return rtn, nil
}

func (r Repo) BillFindAll() ([]Bill, error) {
	var rtn []Bill
	if err := r.bs.C("bills").Find(nil).All(&rtn); err != nil {
		return []Bill{}, err
	}
	return rtn, nil
}

func (r Repo) BillUpdate(selector interface{}, update interface{}) error {
	if err := r.bs.C("bills").Update(selector, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r Repo) BillInsert(bill interface{}) error {
	if err := r.bs.C("bills").Insert(bill); err != nil {
		return err
	}
	return nil
}
