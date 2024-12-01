package main

import (
	"log"
	"monk/database"
	"monk/router"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.InitMongoDB()
	r := mux.NewRouter()
	router.Init(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
