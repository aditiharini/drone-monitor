package api

import (
	"fmt"
	"net/http"
)

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
		Sent     int64
		Received int64
	}
	Saturatr struct {
		Sent     int64
		Received int64
	}
}

func (s *State) HandleDji(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post request from dji")
}

func (s *State) HandleDroneSaturatr(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post request from drone saturatr")
}

func (s *State) HandleServerSaturatr(res http.ResponseWriter, req *http.Request) {
	req.FormValue("acker_packets_sent")
	fmt.Printf(
		"Post request from drone saturatr (acker_packets_sent: %s, acker_packets_received: %s, saturatr_packets_sent: %s, saturatr_packets_received: %s\n",
		req.FormValue("acker_packets_sent"),
		req.FormValue("acker_packets_received"),
		req.FormValue("saturatr_packets_sent"),
		req.FormValue("saturatr_packets_received"),
	)
}

func (s *State) HandleGetState(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Get request from client")
}
