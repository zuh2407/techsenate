package ledger

import (
	"encoding/json"
	"fmt"
	
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	level *leveldb.DB
}

type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

// NewDB opens the LevelDB database at the given path
func NewDB(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{level: db}, nil
}

// Close closes the database
func (db *DB) Close() {
	db.level.Close()
}

// PutTransaction stores a transaction in LevelDB
func (db *DB) PutTransaction(from, to string, amount int) error {
	tx := Transaction{From: from, To: to, Amount: amount}
	key := fmt.Sprintf("tx:%d", time.Now().UnixNano())
	value, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	return db.level.Put([]byte(key), value, nil)
}

// GetAllTransactions retrieves all ledger entries
func (db *DB) GetAllTransactions() []Transaction {
	var results []Transaction
	iter := db.level.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, "tx:") {
			var tx Transaction
			if err := json.Unmarshal(iter.Value(), &tx); err == nil {
				results = append(results, tx)
			}
		}
	}
	iter.Release()
	return results
}
