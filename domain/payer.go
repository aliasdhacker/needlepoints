package domain

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type Payer struct {
	Name    string
	Balance int
}

type PointsTxnRequest struct {
	payer     string
	points    int
	timestamp string
}

type PointsTransaction struct {
	Id        uuid.UUID
	Name      string
	Amount    int
	Timestamp time.Time
	Committed bool
}

var PayerByName map[string]*Payer
var PointsTransactionByName map[string][]*PointsTransaction

// Store the Payer in memory database.
func StorePayer(payer Payer) {
	PayerByName[payer.Name] = &payer
}

// Store the transaction request in memory database.
func StorePointsTransaction(txn PointsTxnRequest) (*PointsTransaction, error) {
	time, err := parseTime(txn.timestamp)

	if err != nil {
		return nil, err
	}

	newT := PointsTransaction{
		Id:        uuid.New(),
		Name:      txn.payer,
		Amount:    txn.points,
		Timestamp: *time,
		Committed: false,
	}

	PointsTransactionByName[newT.Name] = append(PointsTransactionByName[newT.Name], &newT)

	return &newT, nil
}

// Process a transaction against the Payers ledger.
func ProcessTransaction(txn *PointsTransaction) *Payer {
	if txn.Committed {
		return nil
	}

	payer := PayerByName[txn.Name]
	if payer == nil {
		payer = &Payer{
			Name:    txn.Name,
			Balance: 0,
		}
		PayerByName[txn.Name] = payer
	}

	PayerByName[txn.Name].Balance = txn.Amount

	for _, t := range PointsTransactionByName[txn.Name] {
		if txn.Id == t.Id {
			t.Committed = true
		}
	}

	return payer
}

// Utility function to parse time strings in YYYY-MM-DDTHH:MM:SSZ format to golang Time structs
func parseTime(t string) (*time.Time, error) {
	tt, err := time.Parse("2022-01-01T13:00:00Z", t)
	if err != nil {
		log.Fatalf("Failed to parse time: %s", err)
		return nil, errors.New("unable to parse timestamp from provided time value")
	}

	return &tt, nil
}
