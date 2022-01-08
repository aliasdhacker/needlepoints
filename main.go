package main

import (
	"encoding/json"
	"log"
	"needlepoint/domain"
	"net/http"

	"github.com/gorilla/mux"
)

func getAllPayerBalances(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(&domain.PayerByName, "", "\t")
	if err != nil {
		return
	}

	w.Write(output)

	return
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var txn domain.PointsTxnRequest

	json.Unmarshal(body, &txn)

	newT, err := domain.StorePointsTransaction(txn)

	if err != nil {
		errString, _ := json.Marshal(err)
		w.Write(errString)
		return
	}

	p := domain.ProcessTransaction(newT)
	if p != nil {
		output, err := json.MarshalIndent(&p, "", "\t")
		if err != nil {
			return
		}

		w.Write(output)
	}

	return
}

func buy(w http.ResponseWriter, r *http.Request) {
	output, err := json.MarshalIndent(&domain.PayerByName, "", "\t")
	if err != nil {
		return
	}

	w.Write(output)

	return
}

func main() {
	domain.PointsTransactionByName = make(map[string][]*domain.PointsTransaction)
	domain.PayerByName = make(map[string]*domain.Payer)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getAllPayerBalances)
	router.HandleFunc("/getAllPayerBalances", getAllPayerBalances)
	router.HandleFunc("/newTxn", newTransaction)
	router.HandleFunc("/buy", buy)

	log.Fatal(http.ListenAndServe(":8080", router))
}
