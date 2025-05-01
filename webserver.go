package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"go-p2p-ledger/ledger" // replace with correct module import path
)

type Transaction struct {
    From   string `json:"from"`
    To     string `json:"to"`
    Amount int    `json:"amount"`
}

var db *ledger.DB

func main() {
    var err error
    db, err = ledger.NewDB("ledger")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/transfer", transferHandler)
    http.HandleFunc("/ledger", ledgerHandler)

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(filepath.Join("web/", "index.html"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var tx Transaction
    if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Set local sender for simplicity
    tx.From = "USER ACCOUNT"
    err := db.PutTransaction(tx.From, tx.To, tx.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ledgerHandler(w http.ResponseWriter, r *http.Request) {
    entries := db.GetAllTransactions()
    json.NewEncoder(w).Encode(entries)
}
