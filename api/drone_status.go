package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type State struct {
	Drone struct {
		Dji      DjiState      `json:"dji"`
		Saturatr SaturatrState `json:"saturatr"`
	} `json:"drone"`
	Server struct {
		Saturatr SaturatrState `json:"saturatr"`
	} `json:"server"`
	mux sync.Mutex
}

type Point [2]float64

type DjiState struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Location       Point      `json:"location"`
	Battery        float32    `json:"battery"` // fraction of available battery, range [0,1]
	GPS            [3]float64 `json:"gps"`
	XYZ            [3]float64 `json:"xyz"`
	Yaw            float32    `json:"yaw"`
	BatteryVoltage float32    `json:"batteryVoltage"`
	LastSeenTime   int64      `json:"last_seen_time"`
	CruisingSpeed  float64    `json:"cruising_speed"`
}

type SaturatrState struct {
	Acker struct {
		Sent     int64 `json:"sent"`
		Received int64 `json:"received"`
	} `json:"acker"`
	Saturatr struct {
		Sent     int64 `json:"sent"`
		Received int64 `json:"received"`
	} `json:"saturatr"`
}

func (s *State) HandleDji(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post request from dji")
	var dji DjiState
	if err := json.NewDecoder(req.Body).Decode(&dji); err != nil {
		fmt.Printf("Problem handling dji post %v\n", err)
	}
	s.mux.Lock()
	s.Drone.Dji = dji
	s.mux.Unlock()
}

func (s *State) HandleDroneSaturatr(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post request from drone saturatr")
	var saturatr SaturatrState
	if err := json.NewDecoder(req.Body).Decode(&saturatr); err != nil {
		fmt.Printf("Problem handling drone saturatr %v\n", err)
	}
	s.mux.Lock()
	s.Drone.Saturatr = saturatr
	s.mux.Unlock()
}

func (s *State) HandleServerSaturatr(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post request from server saturatr")
	var saturatr SaturatrState
	if err := json.NewDecoder(req.Body).Decode(&saturatr); err != nil {
		fmt.Printf("Problem handling server saturatr %v\n", err)
	}
	s.mux.Lock()
	s.Server.Saturatr = saturatr
	s.mux.Unlock()
}

func (s *State) HandleGetState(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Get request from client")
	(res).Header().Set("Access-Control-Allow-Origin", "*")
	s.mux.Lock()
	json.NewEncoder(res).Encode(s)
	s.mux.Unlock()
}
