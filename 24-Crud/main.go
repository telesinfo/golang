package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// CRUD - CREATE, READ, UPDATE, DELETE

	router := mux.NewRouter()
	router.HandleFunc("/users", server.AddUser).Methods(http.MethodPost)

	fmt.Println("Listen the port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
