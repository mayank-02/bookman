package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mayank-02/bookman/internal/api"
	"github.com/mayank-02/bookman/internal/db"
)

func main() {
	db, err := db.InitDB("./bookman.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	api.RegisterHandlers(router, db)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
