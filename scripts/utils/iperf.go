package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type IperfInfo struct {
	Bandwidth float64 `json:"bandwidth"`
	Unit      string  `json:"unit"`
	Direction string  `json:"direction"`
}

// For realtime print statements
func PostBandwidth(output string, direction string, httpClient http.Client, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "/sec") {
			pieces := strings.Split(strings.TrimSpace(line), "  ")
			fmt.Println(line)
			bandwidthParts := strings.Split(pieces[len(pieces)-1], " ")
			if len(bandwidthParts) < 2 {
				fmt.Println("Couldn't parse due to invalid formatting")
			}
			mbpsStr := strings.TrimSpace(bandwidthParts[0])
			unit := strings.TrimSpace(bandwidthParts[1])
			mbpsFloat, err := strconv.ParseFloat(mbpsStr, 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if direction == "both" {
				if len(pieces) > 6 {
					direction = "upload"
				} else {
					direction = "download"
				}
			}
			body, err := json.Marshal(IperfInfo{Bandwidth: mbpsFloat, Unit: unit, Direction: direction})
			if err != nil {
				fmt.Println(err)
				continue
			}
			res, err := httpClient.Post(endpoint, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(res)
				fmt.Println(err)
			}
		}
	}

}
