package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"needlepoint/domain"
	"net/http"
	"strconv"
)

/*

Self explanatory - get all payer balances

Example:
{ "DANNON": 1100 }
{ "UNILEVER": 200 }
{ "MILLER COORS": 10000 }

*/
func getAllPayerBalances(w http.ResponseWriter, r *http.Request) {
	// Custom output instead of auto marshalling - to match expected results.
	//output, err := json.MarshalIndent(&domain.PayerByName, "", "\t")
	//if err != nil {
	//	return
	//}
	var first bool = true
	var output = "\n{\n"
	for _, payer := range domain.GroupNewTxnsByName(getArrayOfTransactions(domain.PointsTransactionById)) {
		payer := fmt.Sprintf(" \"%s\": %d", payer.Payer, payer.Points)
		if !first {
			payer = ",\n" + payer
		}
		output = fmt.Sprintf("%s%s", output, payer)
		first = false
	}

	w.WriteHeader(200)
	io.WriteString(w, output+"\n}")

	return
}

/*

Submit a points transaction

Response example:
Payer: &{Payer:DANNON Points:1100}

*/
func newTransaction(w http.ResponseWriter, r *http.Request) {
	var txn domain.PointsTxnRequest
	var body []byte
	body, _ = ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &txn)
	newT, err := domain.StorePointsTransactionRequest(txn)

	if err != nil {
		errString, _ := json.Marshal(err)
		w.Write(errString)
		return
	}

	fmt.Fprintf(w, "Transaction: %+v", newT)

	return
}

/*

Spend your points

*/
func spend(w http.ResponseWriter, r *http.Request) {
	var spendRequest domain.SpendRequest
	var output string
	var body []byte
	body, _ = ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &spendRequest)

	responseTransactions, err := domain.CalculateSpend(spendRequest)
	if err == nil {
		output = formatSpendResponse(responseTransactions)
	}

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, output)
	return
}

func main() {
	domain.PointsTransactionById = make(map[uuid.UUID]*domain.PointsTransaction)
	mux := http.NewServeMux()
	mux.HandleFunc("/newTxn", newTransaction)
	mux.HandleFunc("/", getAllPayerBalances)
	mux.HandleFunc("/spend", spend)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// Helper func to convert map to array
func getArrayOfTransactions(txns map[uuid.UUID]*domain.PointsTransaction) []*domain.PointsTransaction {
	var arr []*domain.PointsTransaction
	for _, v := range txns {
		arr = append(arr, v)
	}

	return arr
}

// Helper func - make response match expected response
func formatSpendResponse(transactions []*domain.PointsTransaction) string {
	var output string = "\n[\n"
	var first bool = true

	for _, v := range transactions {
		if !first {
			output += ",\n{"
		} else {
			output += "{"
		}
		output += " \"payer\": \"" + v.Payer + "\", \"points\": " + strconv.Itoa(v.Points) + " }"
		first = false
	}

	output += "\n]\n"

	return output
}
