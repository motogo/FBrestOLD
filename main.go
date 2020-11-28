package main

import (
	"github.com/gorilla/mux"
	"fbrest/Dokumento/apis"
	"log"
	"net/http"
)

func main(){

	router := mux.NewRouter()

	router.HandleFunc("/dokumento/standorte",apis.GetAllLocations).Methods("GET")

	err := http.ListenAndServe(":1234",router)

	if err != nil {
		log.Println(err)
	}

}