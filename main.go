package main

import (
	"github.com/gorilla/mux"

	"github.com/aditiharini/drone-monitor/api"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	state := api.State{}
	router.HandleFunc("/drone/flight", state.HandleDji)
	router.HandleFunc("/drone/saturatr", state.HandleDroneSaturatr)
	router.HandleFunc("/server/saturatr", state.HandleServerSaturatr)
}
