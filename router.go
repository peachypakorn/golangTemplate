package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	AuthMethod  int
}

const Admin = 0
const NoAuth = 0
const LoginNeeded = 0

type Routes []Route

var routes = Routes{
	Route{
		"RecommendedProducts",
		"GET",
		"/users/{uid}/recommendedProducts",
		getRecommendedProducts,
		Admin,
	},
}

var routes3 = Routes{
	Route{
		"StoreCreate",
		"POST",
		"/store",
		StoreCreate,
		Admin,
	},
	Route{
		"StoreUpdate",
		"PUT",
		"/store",
		StoreUpdate,
		Admin,
	},
	Route{
		"StoreDelete",
		"DELETE",
		"/store",
		StoreDelete,
		Admin,
	},
	Route{
		"StoreGetAll",
		"GET",
		"/stores",
		StoreGetAll,
		Admin,
	},
	Route{
		"StoreFind",
		"GET",
		"/store",
		StoreFind,
		Admin,
	},
}

var routes2 = Routes{
	Route{
		"CustomerGetAll",
		"GET",
		"/customers",
		CustomerGetAll,
		0,
	},
	Route{
		"CustomerCreate",
		"POST",
		"/customer/create",
		CustomerCreate,
		0,
	},
	Route{
		"CustomerUpdate",
		"POST",
		"/customer/update",
		CustomerUpdate,
		0,
	},

	Route{
		"AssetGetAll",
		"GET",
		"/assets",
		AssetGetAll,
		0,
	},
	Route{
		"AssetCreate",
		"POST",
		"/asset/create",
		AssetCreate,
		0,
	},
	Route{
		"AssetUpdate",
		"POST",
		"/asset/update",
		AssetUpdate,
		0,
	},

	Route{
		"AliasGetAll",
		"GET",
		"/aliases",
		AliasGetAll,
		0,
	},
	Route{
		"AliasCreate",
		"POST",
		"/alias/create",
		AliasCreate,
		0,
	},
	Route{
		"AliasUpdate",
		"POST",
		"/alias/update",
		AliasUpdate,
		0,
	},

	Route{
		"TransferInquiry",
		"POST",
		"/bank_account/inquiry",
		InquiryBankAccountBalance,
		0,
	},
	Route{
		"TransferInquiry",
		"POST",
		"/Transfer",
		Transfer,
		0,
	},
	Route{
		"ListHistory",
		"POST",
		"/ListHistory",
		ListHistory,
		0,
	},
	Route{
		"ListBankName",
		"GET",
		"/ListBankName",
		ListBankName,
		0,
	},
	Route{
		"CheckDestAccount",
		"POST",
		"/CheckDestAccount",
		CheckDestAccount,
		0,
	},
	Route{
		"AddFavorite",
		"POST",
		"/AddFavorite",
		SetFavorite,
		0,
	},
	Route{
		"ListFavorite",
		"POST",
		"/ListFavorite",
		ListFavorite,
		0,
	},
	Route{
		"CheckDestBillService",
		"POST",
		"/CheckDestBillService",
		CheckDestBillService,
		0,
	},
	Route{
		"BillPayment",
		"POST",
		"/BillPayment",
		BillPayment,
		0,
	},

	Route{
		"BillGetAll",
		"GET",
		"/bills",
		BillGetAll,
		0,
	},
	Route{
		"BillCreate",
		"POST",
		"/bill/create",
		BillCreate,
		0,
	},
	Route{
		"BillUpdate",
		"POST",
		"/bill/update",
		BillUpdate,
		0,
	},
}

const (
	NotAuthen = iota

)

func handleAuthenMethod(handler http.Handler, r Route) http.Handler {
	switch r.AuthMethod {
	case Admin:

	//case NoAuth:
	//
	//case LoginNeeded:
	}
	return handler
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Infof("AuthMethod [%d]: %s\t%s\t%s", route.AuthMethod, route.Name, route.Method, route.Pattern)

		handler := handleAuthenMethod(route.HandlerFunc, route)
		router.
			//PathPrefix("/api/v1").
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
