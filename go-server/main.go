package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)


var db *sql.DB

func main() {
	db = setupDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/create",createShortURL).Methods("POST")
	r.HandleFunc("/api/visit/{shortCode}",visitShortURL).Methods("GET")
	r.HandleFunc("/api/info/{shortCode}",getURLInfo).Methods("GET")

	fmt.
}

