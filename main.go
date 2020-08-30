package main

import (
	"net/http"

	"github.com/aditiharini/drone-monitor/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	state := api.State{}
	router.HandleFunc("/drone/flight", state.HandleDji).Methods("POST")
	router.HandleFunc("/drone/saturatr", state.HandleDroneSaturatr).Methods("POST")
	router.HandleFunc("/server/saturatr", state.HandleServerSaturatr).Methods("POST")
	router.HandleFunc("/state", state.HandleGetState).Methods("GET")
	err := http.ListenAndServe(":10000", router)
	if err != nil {
		panic(err)
	}
}
