package trace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PingInfo struct {
	Latency string `json:"latency"`
}

func PostLatency(output string, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			pieces := strings.Split(line, "time=")
			timeStr := pieces[1]
			ms := strings.Split(timeStr, " ")[0]
			body, err := json.Marshal(PingInfo{ms})
			if err != nil {
				panic(err)
			}
			res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(res)
				fmt.Println(err)
			}
		}
	}
}
