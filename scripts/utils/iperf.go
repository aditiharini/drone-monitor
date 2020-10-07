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
}

// For realtime print statements
func PostBandwidth(output string, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "sec") {
			pieces := strings.Split(line, "  ")
			fmt.Println(line)
			fmt.Println(pieces)
			mbpsStr := strings.TrimSpace(strings.Split(pieces[5], " ")[0])
			mbpsFloat, err := strconv.ParseFloat(mbpsStr, 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Iperf info: ", mbpsFloat)
			body, err := json.Marshal(IperfInfo{Bandwidth: mbpsFloat})
			if err != nil {
				fmt.Println(err)
				continue
			}
			res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(res)
				fmt.Println(err)
			}
		}
	}

}
