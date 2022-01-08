package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"needlepoint/domain"
	"net/http"
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
	var output = "\n"
	for _, payer := range domain.PayerByName {
		payer := fmt.Sprintf("{ \"%s\": %d }\n", payer.Payer, payer.Points)
		output = fmt.Sprintf("%s%s", output, payer)
	}

	w.WriteHeader(200)
	io.WriteString(w, output)

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
	newT, err := domain.StorePointsTransaction(txn)

	if err != nil {
		//errString, _ := json.Marshal(err)

		return
	}

	p := domain.ProcessTransaction(newT)
	if p != nil {
		fmt.Fprintf(w, "Payer: %+v", p)
	}

	return
}

/*

Spend your points

 */
func spend(w http.ResponseWriter, r *http.Request) {
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
	mux := http.NewServeMux()
	mux.HandleFunc("/newTxn", newTransaction)
	mux.HandleFunc("/", getAllPayerBalances)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
