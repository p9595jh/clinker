package repository

import (
	"context"
	"encoding/json"
	"layer/common/util"
	"layer/internal/infrastructure/database/dbconn"
	"layer/internal/infrastructure/database/entity"
	"math/big"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestRepository(t *testing.T) {
	db, _ := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   "127.0.0.1",
		Port:   27017,
		Schema: "layer",
	})
	now := time.Now()

	txnRepository := NewTxnRepository(db.Collection("txns"), util.Transformer())
	id, err := txnRepository.Insert(&entity.Txn{
		Hash: util.RandHex(64),
		// Timestamp: time.Now(), //.Format(util.DateFormat),
		Timestamp: &now,
		GasPrice:  big.NewInt(1000000000),
	})
	t.Log(id, err)

	data, _ := txnRepository.Find(nil)
	t.Log(data)
}

func TestConsole(t *testing.T) {
	db, _ := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   "127.0.0.1",
		Port:   27017,
		Schema: "layer",
	})

	res, err := db.Collection("txns").Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	// var data []bson.M
	var data []entity.Txn
	err = res.All(context.Background(), &data)
	t.Log(err)
	b, _ := json.MarshalIndent(&data, "", "  ")
	t.Log(string(b))
	t.Log(data)
}

func TestFind(t *testing.T) {
	db, _ := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   "127.0.0.1",
		Port:   27017,
		Schema: "layer",
	})

	txnRepository := NewTxnRepository(db.Collection("cas"), util.Transformer())
	txns, _ := txnRepository.Find(bson.M{
		// "timestamp": 1670225090,
		"timestamp": bson.M{
			"$gt": 1670226064,
		},
	})
	t.Log(txns)
}

func TestId(t *testing.T) {
	db, _ := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   "127.0.0.1",
		Port:   27017,
		Schema: "layer",
	})

	txnRepository := NewTxnRepository(db.Collection("txns"), util.Transformer())
	txns, err := txnRepository.Find(nil)
	t.Log(err)
	t.Log(txns)
	b, _ := json.MarshalIndent(&txns, "", "  ")
	t.Log(string(b))
}

func TestCa(t *testing.T) {
	db, _ := dbconn.ConnectDB(&dbconn.DatabaseConfig{
		Host:   "127.0.0.1",
		Port:   27017,
		Schema: "layer",
	})

	caRepository := NewCaRepository(db.Collection("cas"), util.Transformer())
	cas, err := caRepository.Find(nil)
	t.Log(err)
	t.Log(cas)
	b, _ := json.MarshalIndent(&cas, "", "  ")
	t.Log(string(b))
}
