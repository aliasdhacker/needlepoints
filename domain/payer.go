package domain

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

// Definition of a payer
type Payer struct {
	Payer  string
	Points int
}

// Request wrapper for a points transaction
type PointsTxnRequest struct {
	Payer     string `json:"payer"`
	Points    int    `json:"points"`
	Timestamp string `json:"timestamp"`
}

// Request wrapper for spending points
type SpendRequest struct {
	Points int `json:"points"`
}

// Definition of a points transaction
type PointsTransaction struct {
	Id        uuid.UUID
	Payer     string
	Points    int
	Timestamp time.Time
	Committed bool
}

// In memory storage for payers and points transactions
var PayerByName map[string]*Payer
var PointsTransactionByName map[string][]*PointsTransaction

// Store the Payer in memory database.
func StorePayer(payer Payer) {
	PayerByName[payer.Payer] = &payer
}

// Store the transaction request in memory database (after converting it from the request wrapper and parsing time.)
func StorePointsTransaction(txn PointsTxnRequest) (*PointsTransaction, error) {
	time, err := parseTime(txn.Timestamp)

	if err != nil {
		return nil, err
	}

	newT := PointsTransaction{
		Id:        uuid.New(),
		Payer:     txn.Payer,
		Points:    txn.Points,
		Timestamp: *time,
		Committed: false,
	}

	PointsTransactionByName[newT.Payer] = append(PointsTransactionByName[newT.Payer], &newT)

	return &newT, nil
}

// Process a transaction against the Payers ledger.
func ProcessTransaction(txn *PointsTransaction) *Payer {
	if txn.Committed {
		return nil
	}

	payer := PayerByName[txn.Payer]
	if payer == nil {
		payer = &Payer{
			Payer:  txn.Payer,
			Points: 0,
		}
		PayerByName[txn.Payer] = payer
	}

	PayerByName[txn.Payer].Points += txn.Points

	for _, t := range PointsTransactionByName[txn.Payer] {
		if txn.Id == t.Id {
			t.Committed = true
		}
	}

	return payer
}

// Utility function to parse time strings in YYYY-MM-DDTHH:MM:SSZ format to golang Time structs
func parseTime(t string) (*time.Time, error) {
	tt, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		log.Println("Failed to parse time: %s", err)
		return nil, errors.New("unable to parse timestamp from provided time value")
	}

	return &tt, nil
}
