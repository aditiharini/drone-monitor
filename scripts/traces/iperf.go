package trace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type IperfInfo struct {
	Bandwidth string `json:"bandwidth"`
}

// For realtime print statements
func PostBandwidth(output string, endpoint string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "sec") {
			pieces := strings.Split(line, " ")
			bandwidth := pieces[len(pieces)-2]
			body, err := json.Marshal(IperfInfo{Bandwidth: bandwidth})
			res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(res)
				fmt.Println(err)
			}
		}
	}

}
