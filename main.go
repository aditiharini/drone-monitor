package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aditiharini/drone-monitor/api"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	logfile, err := os.Create(fmt.Sprintf("srv-%d.log", time.Now().Unix()))
	logWriter := bufio.NewWriter(logfile)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.StampMicro,
	})
	log.SetOutput(logWriter)
	router := mux.NewRouter().StrictSlash(true)
	state := api.State{}
	go func() {
		// Allow some time for updated information to arrive
		// If it doesn't clear unupdated state.
		for {
			state.ClearUnupdatedState()
			time.Sleep(2 * time.Second)
		}
	}()
	router.HandleFunc("/drone/flight", state.HandleDji).Methods("POST")
	router.HandleFunc("/drone/saturatr", state.HandleDroneSaturatr).Methods("POST")
	router.HandleFunc("/drone/signal", state.HandleDroneSignal).Methods("POST")
	router.HandleFunc("/drone/iperf", state.HandleDroneIperf).Methods("POST")
	router.HandleFunc("/drone/ping", state.HandleDronePing).Methods("POST")
	router.HandleFunc("/server/saturatr", state.HandleServerSaturatr).Methods("POST")
	router.HandleFunc("/server/iperf", state.HandleServerIperf).Methods("POST")
	router.HandleFunc("/state", state.HandleGetState).Methods("GET")
	err = http.ListenAndServe(":10000", router)
	if err != nil {
		panic(err)
	}
}
