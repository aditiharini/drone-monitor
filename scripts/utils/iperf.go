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
}

// For realtime print statements
func PostBandwidth(output string, httpClient http.Client, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "sec") {
			pieces := strings.Split(line, "  ")
			fmt.Println(line)
			fmt.Println(pieces)
			bandwidthParts := strings.Split(pieces[5], " ")
			mbpsStr := strings.TrimSpace(bandwidthParts[0])
			unit := strings.TrimSpace(bandwidthParts[1])
			mbpsFloat, err := strconv.ParseFloat(mbpsStr, 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Iperf info: ", mbpsFloat)
			body, err := json.Marshal(IperfInfo{Bandwidth: mbpsFloat, Unit: unit})
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
