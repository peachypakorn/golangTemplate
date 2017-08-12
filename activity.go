package main

import (
	"encoding/json"
	"net/http"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func InquiryBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	log.Debugln("InquiryBankAccountBalance")
	var req struct {
		BeaconID string `json:"beaconId"`
		BankName string `json:"bankName"`
		Alias    string `json:"sourceAccAliasName"`
	}

	type res struct {
		Code          string  `json:"responseCode"`
		Description   string  `json:"responseDescription"`
		AccountName   string  `json:"sourceAccName"`
		BankName      string  `json:"bankName"`
		AccountNumber string  `json:"sourceAccNo"`
		Balance       float64 `json:"balanceAmount"`
		Currency      string  `json:"currency"`
		AccountStatus string  `json:"accStatus"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	alias, err := db.AliasFind(bson.M{"owner_id": bson.ObjectIdHex(req.BeaconID), "alias": req.Alias})
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	asset, err := db.AssetFindByID(alias.AssetID.Hex())
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	JSONResponse(w, http.StatusFound, res{
		Code:          "00",
		Description:   "success",
		AccountName:   asset.BankName,
		BankName:      asset.BankName,
		AccountNumber: asset.AccNumber,
		Balance:       asset.Balance,
		Currency:      asset.Currency,
		AccountStatus: asset.Status,
	})

}

func Transfer(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Transfer")
	type SourceAccount struct {
		BankName      string `json:"bankName"`
		AccountNumber string `json:"sourceAccNo"`
		Alias         string `json:"sourceAccAliasName"`
	}

	var req struct {
		BeaconID                 string        `json:"beaconId"`
		SourceAccount            SourceAccount `json:"account"`
		DestinationBankName      string        `json:"destBankName"`
		DestinationAccountNumber string        `json:"destAccNo"`
		Amount                   float64       `json:"transferAmount"`
		Currency                 string        `json:"currency"`
	}

	type res struct {
		Code          string  `json:"responseCode"`
		Description   string  `json:"responseDescription"`
		SourceBalance float64 `json:"sourceBalanceAmount"`
		Currency      string  `json:"currency"`
		Reference     string  `json:"ref"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	var sAcc Asset
	if req.SourceAccount.AccountNumber != "" && req.SourceAccount.BankName != "" {
		s := bson.M{"bank_name": req.SourceAccount.BankName, "account_number": req.SourceAccount.AccountNumber}
		var err error
		sAcc, err = db.AssetFind(s)
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	} else {
		alias, err := db.AliasFind(bson.M{"owner_id": bson.ObjectIdHex(req.BeaconID), "alias": req.SourceAccount.Alias})
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}

		sAcc, err = db.AssetFindByID(alias.AssetID.Hex())
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	}

	dAcc, err := db.AssetFind(bson.M{"bank_name": req.DestinationBankName, "account_number": req.DestinationAccountNumber})
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	if sAcc.Balance < req.Amount {
		log.Errorln("insufficient balance to do transaction:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: "insufficient balance to do transaction"})
		return
	}

	if (sAcc.Status != "active") || (dAcc.Status != "active") {
		log.Errorln("source or destination account is not active:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: "source or destination account is not active"})
		return
	}

	sAcc.Balance = sAcc.Balance - req.Amount
	dAcc.Balance = dAcc.Balance + req.Amount

	err = db.AssetUpdate(bson.M{"_id": sAcc.ID}, sAcc)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	err = db.AssetUpdate(bson.M{"_id": dAcc.ID}, dAcc)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	t := Transaction{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.ObjectIdHex(req.BeaconID),
		CreatedAt: time.Now(),
		Type:      "transfer",

		SourceBankName: sAcc.BankName,
		SourceAccNo:    sAcc.AccNumber,
		DestBankName:   dAcc.BankName,
		DestAccNo:      dAcc.AccNumber,
		TransferAmount: req.Amount,
		Currency:       req.Currency,
		TransCode:      "TRANSFER",
	}

	if err := db.TransactionInsert(t); err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	JSONResponse(w, http.StatusOK, res{
		Code:          "00",
		Description:   "success",
		Currency:      "THB",
		SourceBalance: sAcc.Balance,
		Reference:     t.ID.Hex(),
	})
}

var TransCodeMap = map[string]string{
	"TRANSFER":    "transfer",
	"TOPUP":       "top up",
	"BILLPAYMENT": "bill payment",
}

type Transaction struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	OwnerID   bson.ObjectId `json:"-" bson:"owner_id"`
	CreatedAt time.Time     `json:"-" bson:"created_at" bson:"created_at"`
	Type      string        `json:"-" bson:"type"`

	// Default
	SourceBankName string `json:"sourceBankName,omitempty" bson:"sourceBankName,omitempty"`
	SourceAccNo    string `json:"sourceAccNo,omitempty" bson:"souceAccNo,omitempty"`
	Date           string `json:"date,omitempty" bson:"-"`
	Time           string `json:"time,omitempty" bson:"-"`
	TransCode      string `json:"transCode,omitempty" bson:"transCode,omitempty"`
	TransCodeDesc  string `json:"transCodeDesc,omitempty" bson:"-"`
	Currency       string `json:"currency,omitempty" bson:"currency,omitempty"`

	// Transfer
	DestBankName   string  `json:"destBankName,omitempty" bson:"destBankName,omitempty"`
	DestAccNo      string  `json:"destAccNo,omitempty" bson:"destAccNo,omitempty"`
	TransferAmount float64 `json:"transferAmount,omitempty" bson:"transferAmount,omitempty"`

	// TopUp
	TopUpService       string  `json:"topupService,omitempty" bson:"topupService,omitempty"`
	Provider           string  `json:"serviceProvider,omitempty" bson:"serviceProvider,omitempty"`
	Ref1               string  `json:"REF1,omitempty" bson:"REF1,omitempty"`
	Ref2               string  `json:"REF2,omitempty" bson:"REF2,omitempty"`
	TopUpAccountNumber string  `json:"topupAccNo,omitempty" bson:"topupAccNo,omitempty"`
	TopUpAmount        float64 `json:"topupAmount,omitempty" bson:"topupAmount,omitempty"`

	// BillPayment
	BillID            bson.ObjectId `json:"-" bson:"bill_id,omitempty"`
	SourceCardNumber  string        `json:"sourceCardNo,omitempty" bson:"sourceCardNo,omitempty"`
	BillType          string        `json:"billType,omitempty" bson:"billType,omitempty"`
	BillService       string        `json:"billService,omitempty" bson:"billService,omitempty"`
	BillAccountNumber string        `json:"billAccNo,omitempty" bson:"billAccNo,omitempty"`
	BillRef2          string        `json:"billREF1,omitempty" bson:"billREF1,omitempty"`
	BillRef1          string        `json:"billREF2,omitempty" bson:"billREF2,omitempty"`
	PaidAmount        float64       `json:"paidAmount,omitempty" bson:"paidAmount,omitempty"`

	// Favorite
	FavoriteAlias      string `json:"favoriteAliasName,omitempty" bson:"favorite,omitempty"`
	SourceAccAliasName string `json:"-" bson:"sourceAccAliasName,omitempty"`
	DestPhoneNo        string `json:"-" bson:"destPhoneNo,omitempty"`
}

type TransactionTransfer struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	OwnerID   bson.ObjectId `json:"-" bson:"owner_id"`
	CreatedAt time.Time     `json:"-" bson:"created_at" bson:"created_at"`
	Type      string        `json:"-" bson:"type"`

	// Default
	SourceBankName string `json:"sourceBankName,omitempty" bson:"sourceBankName,omitempty"`
	SourceAccNo    string `json:"sourceAccNo,omitempty" bson:"souceAccNo,omitempty"`
	Date           string `json:"date,omitempty" bson:"-"`
	Time           string `json:"time,omitempty" bson:"-"`
	TransCode      string `json:"transCode,omitempty" bson:"transCode,omitempty"`
	TransCodeDesc  string `json:"transCodeDesc,omitempty" bson:"-"`
	Currency       string `json:"currency,omitempty" bson:"currency,omitempty"`

	// Transfer
	DestBankName   string  `json:"destBankName,omitempty" bson:"destBankName,omitempty"`
	DestAccNo      string  `json:"destAccNo,omitempty" bson:"destAccNo,omitempty"`
	TransferAmount float64 `json:"transferAmount,omitempty" bson:"transferAmount,omitempty"`

	// Favorite
	FavoriteAlias      string `json:"favoriteAliasName,omitempty" bson:"favorite,omitempty"`
	SourceAccAliasName string `json:"sourceAccAliasName" bson:"sourceAccAliasName,omitempty"`
	DestPhoneNo        string `json:"destPhoneNo" bson:"destPhoneNo,omitempty"`
}

type TransactionTopUp struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	OwnerID   bson.ObjectId `json:"-" bson:"owner_id"`
	CreatedAt time.Time     `json:"-" bson:"created_at" bson:"created_at"`
	Type      string        `json:"-" bson:"type"`

	// Default
	SourceBankName string `json:"sourceBankName,omitempty" bson:"sourceBankName,omitempty"`
	SourceAccNo    string `json:"sourceAccNo,omitempty" bson:"souceAccNo,omitempty"`
	Date           string `json:"date,omitempty" bson:"-"`
	Time           string `json:"time,omitempty" bson:"-"`
	TransCode      string `json:"transCode,omitempty" bson:"transCode,omitempty"`
	TransCodeDesc  string `json:"transCodeDesc,omitempty" bson:"-"`
	Currency       string `json:"currency,omitempty" bson:"currency,omitempty"`

	// TopUp and BillPay
	Provider string `json:"serviceProvider,omitempty" bson:"serviceProvider,omitempty"`
	Ref1     string `json:"REF1,omitempty" bson:"REF1,omitempty"`
	Ref2     string `json:"REF2,omitempty" bson:"REF2,omitempty"`

	// TopUp
	TopUpService       string  `json:"topupService,omitempty" bson:"topupService,omitempty"`
	TopUpAccountNumber string  `json:"topupAccNo,omitempty" bson:"topupAccNo,omitempty"`
	TopUpAmount        float64 `json:"topupAmount,omitempty" bson:"topupAmount,omitempty"`

	// Favorite
	FavoriteAlias      string `json:"favoriteAliasName,omitempty" bson:"favorite,omitempty"`
	SourceAccAliasName string `json:"-" bson:"sourceAccAliasName,omitempty"`
	DestPhoneNo        string `json:"-" bson:"destPhoneNo,omitempty"`
}

type TransactionBillPayment struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	OwnerID   bson.ObjectId `json:"-" bson:"owner_id"`
	CreatedAt time.Time     `json:"-" bson:"created_at" bson:"created_at"`
	Type      string        `json:"-" bson:"type"`

	// Default
	SourceBankName string `json:"sourceBankName,omitempty" bson:"sourceBankName,omitempty"`
	SourceAccNo    string `json:"sourceAccNo,omitempty" bson:"souceAccNo,omitempty"`
	Date           string `json:"date,omitempty" bson:"-"`
	Time           string `json:"time,omitempty" bson:"-"`
	TransCode      string `json:"transCode,omitempty" bson:"transCode,omitempty"`
	TransCodeDesc  string `json:"transCodeDesc,omitempty" bson:"-"`
	Currency       string `json:"currency,omitempty" bson:"currency,omitempty"`

	// TopUp and BillPay
	Provider string `json:"serviceProvider,omitempty" bson:"serviceProvider,omitempty"`
	Ref1     string `json:"REF1,omitempty" bson:"REF1,omitempty"`
	Ref2     string `json:"REF2,omitempty" bson:"REF2,omitempty"`

	// BillPayment
	SourceCardNumber  string  `json:"sourceCardNo,omitempty" bson:"sourceCardNo,omitempty"`
	BillType          string  `json:"billType,omitempty" bson:"billType,omitempty"`
	BillService       string  `json:"billService,omitempty" bson:"billService,omitempty"`
	BillAccountNumber string  `json:"billAccNo,omitempty" bson:"billAccNo,omitempty"`
	BillRef2          string  `json:"billREF1,omitempty" bson:"billREF1,omitempty"`
	BillRef1          string  `json:"billREF2,omitempty" bson:"billREF2,omitempty"`
	PaidAmount        float64 `json:"paidAmount,omitempty" bson:"paidAmount,omitempty"`

	// Favorite
	FavoriteAlias      string `json:"favoriteAliasName,omitempty" bson:"favorite,omitempty"`
	SourceAccAliasName string `json:"-" bson:"sourceAccAliasName,omitempty"`
}

func ListHistory(w http.ResponseWriter, r *http.Request) {
	log.Debugln("ListHistory")
	var req struct {
		BeaconID        string `json:"beaconId"`
		TransactionType string `json:"type"`
		MaxRecord       int    `json:"maxRecord"`
	}

	type res struct {
		Code                    string        `json:"responseCode"`
		Description             string        `json:"responseDescription"`
		TransferTransactions    []Transaction `json:"transfer,omitempty"`
		TopUpTransactions       []Transaction `json:"topup,omitempty"`
		BillPaymentTransactions []Transaction `json:"billPayment,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	h, err := db.TransactionFindByType(req.BeaconID, req.TransactionType, req.MaxRecord)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	var rsp res
	switch req.TransactionType {
	case "transfer":
		rsp.TransferTransactions = h
		for i := range rsp.TransferTransactions {
			rsp.TransferTransactions[i].Date = rsp.TransferTransactions[i].CreatedAt.Format("20060102")
			rsp.TransferTransactions[i].Time = rsp.TransferTransactions[i].CreatedAt.Format("15:04:05")
			rsp.TransferTransactions[i].TransCodeDesc = TransCodeMap[rsp.TransferTransactions[i].TransCode]
		}
	case "topup":
		rsp.TopUpTransactions = h
		for i := range rsp.TopUpTransactions {
			rsp.TopUpTransactions[i].Date = rsp.TopUpTransactions[i].CreatedAt.Format("20060102")
			rsp.TopUpTransactions[i].Time = rsp.TopUpTransactions[i].CreatedAt.Format("15:04:05")
			rsp.TopUpTransactions[i].TransCodeDesc = TransCodeMap[rsp.TopUpTransactions[i].TransCode]
		}
	case "billPayment":
		rsp.BillPaymentTransactions = h
		for i := range rsp.BillPaymentTransactions {
			rsp.BillPaymentTransactions[i].Date = rsp.BillPaymentTransactions[i].CreatedAt.Format("20060102")
			rsp.BillPaymentTransactions[i].Time = rsp.BillPaymentTransactions[i].CreatedAt.Format("15:04:05")
			rsp.BillPaymentTransactions[i].TransCodeDesc = TransCodeMap[rsp.BillPaymentTransactions[i].TransCode]
		}
	}
	rsp.Code = "00"
	rsp.Description = "success"

	JSONResponse(w, http.StatusFound, rsp)
}

func SetFavorite(w http.ResponseWriter, r *http.Request) {
	log.Debugln("AddFavorite")
	var req struct {
		BeaconID      string `json:"beaconId"`
		Reference     string `json:"ref"`
		Favorite      string `json:"favorite"`
		FavoriteAlias string `json:"favoriteAliasName"`
	}

	type res struct {
		Code        string `json:"responseCode"`
		Description string `json:"responseDescription"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if req.Favorite == "Y" {
		err := db.TransactionAddFavorite(req.Reference, req.FavoriteAlias)
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	} else if req.Favorite == "N" {
		err := db.TransactionUnsetField(req.Reference, "favorite")
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	} else {
		log.Errorln("favorite field must be Y or N, but is:", req.Favorite)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: "favorite field must be Y or N"})
		return
	}

	JSONResponse(w, http.StatusOK, res{Code: "00", Description: "success"})
	return
}

func ListFavorite(w http.ResponseWriter, r *http.Request) {
	log.Debugln("ListFavorites")
	var req struct {
		BeaconID        string `json:"beaconId"`
		TransactionType string `json:"type"`
	}

	type res struct {
		Code                    string                   `json:"responseCode"`
		Description             string                   `json:"responseDescription"`
		TransferTransactions    []TransactionTransfer    `json:"transfer,omitempty"`
		TopUpTransactions       []TransactionTopUp       `json:"topup,omitempty"`
		BillPaymentTransactions []TransactionBillPayment `json:"billPayment,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	favorites, err := db.TransactionFindFavorites(req.BeaconID, req.TransactionType)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	var rsp res
	switch req.TransactionType {
	case "transfer":
		rsp.TransferTransactions = []TransactionTransfer{}
		for _, f := range favorites {
			t := TransactionTransfer{
				ID:                 f.ID,
				OwnerID:            f.OwnerID,
				SourceAccAliasName: f.SourceAccAliasName,
				SourceBankName:     f.SourceBankName,
				SourceAccNo:        f.SourceAccNo,
				Date:               f.CreatedAt.Format("20060102"),
				Time:               f.CreatedAt.Format("15:04:05"),
				TransCode:          f.TransCode,
				TransCodeDesc:      TransCodeMap[f.TransCode],
				DestBankName:       f.DestBankName,
				DestAccNo:          f.DestAccNo,
				TransferAmount:     f.TransferAmount,
				Currency:           f.Currency,
				DestPhoneNo:        f.DestPhoneNo,
				FavoriteAlias:      f.FavoriteAlias,
			}
			rsp.TransferTransactions = append(rsp.TransferTransactions, t)
		}
	case "topup":
		rsp.TopUpTransactions = []TransactionTopUp{}
		for _, f := range favorites {
			t := TransactionTopUp{
				ID:                 f.ID,
				OwnerID:            f.OwnerID,
				CreatedAt:          f.CreatedAt,
				Type:               f.Type,
				SourceBankName:     f.SourceBankName,
				SourceAccNo:        f.SourceAccNo,
				Date:               f.CreatedAt.Format("20060102"),
				Time:               f.CreatedAt.Format("15:04:05"),
				TransCode:          f.TransCode,
				TransCodeDesc:      TransCodeMap[f.TransCode],
				Currency:           f.Currency,
				TopUpService:       f.TopUpService,
				Provider:           f.Provider,
				Ref1:               f.Ref1,
				Ref2:               f.Ref2,
				TopUpAccountNumber: f.TopUpAccountNumber,
				TopUpAmount:        f.TopUpAmount,
				FavoriteAlias:      f.FavoriteAlias,
				SourceAccAliasName: f.SourceAccAliasName,
				DestPhoneNo:        f.DestPhoneNo,
			}
			rsp.TopUpTransactions = append(rsp.TopUpTransactions, t)
		}
	case "billPayment":
		rsp.BillPaymentTransactions = []TransactionBillPayment{}
		for _, f := range favorites {
			t := TransactionBillPayment{
				ID:                 f.ID,
				OwnerID:            f.OwnerID,
				CreatedAt:          f.CreatedAt,
				Type:               f.Type,
				SourceBankName:     f.SourceBankName,
				SourceAccNo:        f.SourceAccNo,
				Date:               f.CreatedAt.Format("20060102"),
				Time:               f.CreatedAt.Format("15:04:05"),
				TransCode:          f.TransCode,
				TransCodeDesc:      TransCodeMap[f.TransCode],
				Currency:           f.Currency,
				Provider:           f.Provider,
				Ref1:               f.Ref1,
				Ref2:               f.Ref2,
				SourceCardNumber:   f.SourceCardNumber,
				BillType:           f.BillType,
				BillService:        f.BillService,
				BillAccountNumber:  f.BillAccountNumber,
				BillRef2:           f.BillRef2,
				BillRef1:           f.BillRef1,
				PaidAmount:         f.PaidAmount,
				FavoriteAlias:      f.FavoriteAlias,
				SourceAccAliasName: f.SourceAccAliasName,
			}
			rsp.BillPaymentTransactions = append(rsp.BillPaymentTransactions, t)
		}
	}

	rsp.Code = "00"
	rsp.Description = "success"

	JSONResponse(w, http.StatusFound, rsp)
}

func CheckDestBillService(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BillType        string `json:"billType,omitempty"`
		BillService     string `json:"billService,omitempty"`
		ServiceProvider string `json:"serviceProvider,omitempty"`
		Ref1            string `json:"REF1,omitempty"`
		Ref2            string `json:"REF2,omitempty"`
		BillAccNo       string `json:"billAccNo,omitempty"`
		BillRef1        string `json:"billREF1,omitempty"`
		BillRef2        string `json:"billREF2,omitempty"`
	}

	type res struct {
		Code              string  `json:"responseCode"`
		Description       string  `json:"responseDescription"`
		DestAccountName   string  `json:"destAccName"`
		BillType          string  `json:"billType"`
		BillService       string  `json:"billService"`
		ServiceProvider   string  `json:"serviceProvider"`
		Ref1              string  `json:"REF1"`
		Ref2              string  `json:"REF2"`
		billAccountNumber string  `json:"billAccNo"`
		BillRef1          string  `json:"billREF1"`
		BillRef2          string  `json:"billREF2"`
		DueDate           string  `json:"dueDate"`
		BilledAmount      float64 `json:"billedAmount"`
		MinAmount         float64 `json:"minAmount"`
		Currency          string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	var selector bson.M
	if req.BillRef1 != "" {
		selector = bson.M{
			"billREF1": req.BillRef1,
		}
	} else {
		selector = bson.M{
			"billType":        req.BillType,
			"billService":     req.BillService,
			"serviceProvider": req.ServiceProvider,
			"billAccNo":       req.BillAccNo,
		}
	}

	b, err := db.BillFind(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	JSONResponse(w, http.StatusFound, res{
		Code:              "00",
		Description:       "success",
		BillType:          b.BillType,
		BillService:       b.BillService,
		ServiceProvider:   b.ServiceProvider,
		Ref1:              b.Ref1,
		Ref2:              b.Ref2,
		billAccountNumber: b.BillAccNo,
		BillRef1:          b.BillRef1,
		BillRef2:          b.BillRef2,
		DueDate:           b.DueDate.Format("20060102"),
		BilledAmount:      b.BilledAmount,
		MinAmount:         b.MinAmount,
		Currency:          b.Currency,
	})
}

func BillPayment(w http.ResponseWriter, r *http.Request) {
	log.Debugln("BillPayment")

	type SourceAccount struct {
		BankName      string `json:"bankName"`
		AccountNumber string `json:"sourceAccNo"`
		Alias         string `json:"sourceAccAliasName"`
	}

	var req struct {
		BeaconID        string        `json:"beaconId"`
		SourceAccount   SourceAccount `json:"account"`
		BillType        string        `json:"billType"`
		BillService     string        `json:"billService"`
		ServiceProvider string        `json:"serviceProvider"`
		REF1            string        `json:"REF1"`
		REF2            string        `json:"REF2"`
		BillAccNo       string        `json:"billAccNo"`
		BillRef1        string        `json:"billREF1"`
		BillREF2        string        `json:"billREF2"`
		PaidAmount      float64       `json:"paidAmount"`
		Currency        string        `json:"currency"`
	}

	type res struct {
		Code          string  `json:"responseCode"`
		Description   string  `json:"responseDescription"`
		SourceBalance float64 `json:"sourceBalanceAmount"`
		Currency      string  `json:"currency"`
		Ref           string  `json:"ref"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	var sAcc Asset
	if req.SourceAccount.AccountNumber != "" && req.SourceAccount.BankName != "" {
		s := bson.M{"bank_name": req.SourceAccount.BankName, "account_number": req.SourceAccount.AccountNumber}
		var err error
		sAcc, err = db.AssetFind(s)
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	} else {
		alias, err := db.AliasFind(bson.M{"owner_id": bson.ObjectIdHex(req.BeaconID), "alias": req.SourceAccount.Alias})
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}

		sAcc, err = db.AssetFindByID(alias.AssetID.Hex())
		if err != nil {
			log.Errorln("mongo error:", err)
			JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
			return
		}
	}

	var selector bson.M
	if req.BillRef1 != "" {
		selector = bson.M{
			"billREF1": req.BillRef1,
		}
	} else {
		selector = bson.M{
			"billType":        req.BillType,
			"billService":     req.BillService,
			"serviceProvider": req.ServiceProvider,
			"billAccNo":       req.BillAccNo,
		}
	}

	b, err := db.BillFind(selector)
	if err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	if sAcc.Balance < req.PaidAmount {
		msg := fmt.Sprintf("insufficient balance, balance: %f, paid: %f", sAcc.Balance, req.PaidAmount)
		log.Errorf(msg)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: msg})
		return
	}

	b.PaidAmount = b.PaidAmount + req.PaidAmount
	db.BillUpdate(bson.M{"_id": b.ID}, bson.M{
		"paidAmount": b.PaidAmount,
	})

	t := Transaction{
		ID:        bson.NewObjectId(),
		OwnerID:   bson.ObjectIdHex(req.BeaconID),
		CreatedAt: time.Now(),
		Type:      "billPayment",
		BillID:    b.ID,

		SourceBankName:    sAcc.BankName,
		SourceAccNo:       sAcc.AccNumber,
		SourceCardNumber:  "",
		BillType:          b.BillType,
		BillService:       b.BillService,
		BillAccountNumber: "",
		BillRef1:          b.Ref1,
		BillRef2:          b.Ref2,
		PaidAmount:        req.PaidAmount,
		Currency:          b.Currency,
		TransCode:         "BILLPAYMENT",
	}

	if err := db.TransactionInsert(t); err != nil {
		log.Errorln("mongo error:", err)
		JSONResponse(w, http.StatusBadRequest, res{Code: "10", Description: err.Error()})
		return
	}

	JSONResponse(w, http.StatusOK, res{
		Code:          "00",
		Description:   "success",
		Currency:      "THB",
		SourceBalance: sAcc.Balance,
		Ref:           t.ID.Hex(),
	})
}

func BillCreate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("BillCreate")

	var req Bill
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Debugln(req)

	db := RP.Request(r)
	if err := db.BillInsert(req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func BillUpdate(w http.ResponseWriter, r *http.Request) {
	log.Debugln("BillUpdate")

	var req Bill
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorln("decode json error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db := RP.Request(r)
	if err := db.BillUpdate(bson.M{"_id": req.ID}, req); err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func BillGetAll(w http.ResponseWriter, r *http.Request) {
	log.Debugln("BillGetAll")

	db := RP.Request(r)
	a, err := db.BillFindAll()
	if err != nil {
		log.Errorln("mongo error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JSONResponse(w, http.StatusFound, a)
}
