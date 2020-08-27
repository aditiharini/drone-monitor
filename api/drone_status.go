package api

import "net/http"

type State struct {
	Drone struct {
		Dji      DjiState
		Saturatr SaturatrState
	}
	Server struct {
		Saturatr SaturatrState
	}
}
type Point [2]float64

type DjiState struct {
	Location       Point      `json:"location"`
	Battery        float32    `json:"battery"` // fraction of available battery, range [0,1]
	GPS            [3]float64 `json:"gps"`
	Yaw            float32    `json:"yaw"`
	BatteryVoltage float32    `json:"batteryVoltage"`
	LastSeenTime   int64      `json:"last_seen_time"`
	CruisingSpeed  float64    `json:"cruising_speed"`
}

type SaturatrState struct {
	Acker struct {
		Sent     int64
		Received int64
	}
	Saturatr struct {
		Sent     int64
		Received int64
	}
}

func (s *State) HandleDji(res http.ResponseWriter, req *http.Request) {

}

func (s *State) HandleDroneSaturatr(res http.ResponseWriter, req *http.Request) {

}

func (s *State) HandleServerSaturatr(res http.ResponseWriter, req *http.Request) {

}