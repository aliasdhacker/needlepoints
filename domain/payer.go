package domain

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

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

// In memory storage - points transactions
var PointsTransactionById map[uuid.UUID]*PointsTransaction

// Store the transaction request in memory database (after converting it from the request wrapper and parsing time.)
func StorePointsTransactionRequest(txn PointsTxnRequest) (*PointsTransaction, error) {
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

	StorePointsTransaction(&newT)

	return &newT, nil
}

// Programmatic add of points transaction
func StorePointsTransaction(txn *PointsTransaction) *PointsTransaction {
	PointsTransactionById[txn.Id] = txn

	return txn
}

func GetOldestPoints(pointsTxns map[uuid.UUID]*PointsTransaction) *PointsTransaction {
	var oldest *PointsTransaction

	for _, i := range pointsTxns {
		if oldest == nil || i.Timestamp.Before(oldest.Timestamp) {
			oldest = i
		}
	}

	return oldest
}

func CalculateSpend(request SpendRequest) ([]*PointsTransaction, error) {
	var totalSpent int
	var newTxns []*PointsTransaction
	var allTransactions = make(map[uuid.UUID]*PointsTransaction)

	for k, v := range PointsTransactionById {
		allTransactions[k] = v
	}

	for request.Points > totalSpent {
		txn := GetOldestPoints(allTransactions)
		if txn == nil {
			return nil, errors.New("Spent more points than available - only spending " + strconv.Itoa(totalSpent))
		}
		points := 0
		if txn.Points > request.Points {
			points = request.Points
			request.Points = 0
		} else {
			points = txn.Points
			request.Points -= txn.Points
		}

		newT := PointsTransaction{
			Id:        uuid.New(),
			Payer:     txn.Payer,
			Points:    points * -1,
			Timestamp: time.Now(),
			Committed: true,
		}

		StorePointsTransaction(&newT)
		totalSpent += points
		delete(allTransactions, txn.Id)
		newTxns = append(newTxns, &newT)
	}

	return GroupNewTxnsByName(newTxns), nil

}

func GroupNewTxnsByName(txns []*PointsTransaction) []*PointsTransaction {
	var txnsGroupedByName []*PointsTransaction

	for _, v := range txns {
		if contains(txnsGroupedByName, v) {
		} else {
			copy := &PointsTransaction{
				Id:        v.Id,
				Payer:     v.Payer,
				Points:    v.Points,
				Timestamp: v.Timestamp,
				Committed: v.Committed,
			}
			txnsGroupedByName = append(txnsGroupedByName, copy)
		}
	}

	return txnsGroupedByName
}

// Test an array to see if it contains a value, if it does, add points
func contains(container []*PointsTransaction, txn *PointsTransaction) bool {
	for _, v := range container {
		if v.Payer == txn.Payer {
			v.Points += txn.Points
			return true
		}
	}

	return false
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
