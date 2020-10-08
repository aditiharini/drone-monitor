package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PingInfo struct {
	Latency float64 `json:"latency"`
}

func PostLatency(output string, client http.Client, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			pieces := strings.Split(line, "time=")
			timeStr := pieces[1]
			msStr := strings.TrimSpace(strings.Split(timeStr, " ")[0])
			msFloat, err := strconv.ParseFloat(msStr, 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			body, err := json.Marshal(PingInfo{msFloat})
			if err != nil {
				fmt.Println(err)
				continue
			}
			res, err := client.Post(endpoint, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(res)
				fmt.Println(err)
			}
		}
	}
}
