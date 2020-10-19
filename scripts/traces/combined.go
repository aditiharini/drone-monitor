package trace

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aditiharini/drone-monitor/api"
)

type CombinedTrace struct {
	OutputDir string
	Filepath  string
}

type MicroTime struct {
	time.Time
}

type LogLine struct {
	Time  MicroTime `json:"time"`
	State api.State `json:"state"`
}

func (mt *MicroTime) UnmarshalJSON(buf []byte) error {
	parsedTime, err := time.ParseInLocation(time.StampMicro, strings.Trim(string(buf), `"`), time.Now().Location())
	if err != nil {
		return err
	}
	mt.Time = parsedTime
	return nil
}

func (ll *LogLine) ToCsvRow(baseTime time.Time) []string {
	offsetTime := ll.Time.Sub(baseTime)
	timeStr := fmt.Sprintf("%f", offsetTime.Seconds())
	latitude := fmt.Sprintf("%f", ll.State.Drone.Dji.GPS[0])
	if ll.State.Drone.Dji.GPS[0] == 0 {
		latitude = "NA"
	}
	longitude := fmt.Sprintf("%f", ll.State.Drone.Dji.GPS[1])
	if ll.State.Drone.Dji.GPS[1] == 0 {
		longitude = "NA"
	}
	altitude := fmt.Sprintf("%f", ll.State.Drone.Dji.GPS[2])
	if ll.State.Drone.Dji.GPS[1] == 0 && ll.State.Drone.Dji.GPS[0] == 0 {
		altitude = "NA"
	}
	rsrp := ll.State.Drone.Signal.Rsrp
	if rsrp == "-1" || rsrp == "" {
		rsrp = "NA"
	} else {
		rsrp = rsrp[:len(rsrp)-3]
	}
	rsrq := ll.State.Drone.Signal.Rsrq
	if rsrq == "-1" || rsrq == "" {
		rsrq = "NA"
	} else {
		rsrq = rsrq[:len(rsrq)-2]
	}
	rssi := ll.State.Drone.Signal.Rssi
	if rssi == "-1" || rssi == "" {
		rssi = "NA"
	}
	sinr := ll.State.Drone.Signal.Sinr
	if sinr == "-1" || sinr == "" {
		sinr = "NA"
	}
	cellId := ll.State.Drone.Signal.CellId
	if cellId == "-1" || cellId == "" {
		cellId = "NA"
	}
	uplinkBandwidth := fmt.Sprintf("%f", ll.State.Server.Iperf.Bandwidth)
	if ll.State.Server.Iperf.Bandwidth == -1 {
		uplinkBandwidth = "NA"
	} else if ll.State.Server.Iperf.Unit == "Kbits/sec" {
		uplinkBandwidth = fmt.Sprintf("%f", ll.State.Server.Iperf.Bandwidth/1000)
	}
	downlinkBandwidth := fmt.Sprintf("%f", ll.State.Drone.Iperf.Bandwidth)
	if ll.State.Drone.Iperf.Bandwidth == -1 {
		downlinkBandwidth = "NA"
	}
	latency := fmt.Sprintf("%f", ll.State.Drone.Ping.Latency)
	if ll.State.Drone.Ping.Latency == -1 {
		latency = "NA"
	}
	// TODO(aditi): change bw and latency to be parsed as strings
	return []string{timeStr, latitude, longitude, altitude, rsrp, rsrq, rssi, sinr, cellId, uplinkBandwidth, downlinkBandwidth, latency}
}

func (ct *CombinedTrace) CsvHeaders() []string {
	return []string{"time", "latitude", "longitude", "altitude", "rsrp", "rsrq", "rssi", "sinr", "cellId", "uplinkBw", "downlinkBw", "latency"}
}

func (ct *CombinedTrace) Filename() string {
	sections := strings.Split(ct.Filepath, "/")
	nameParts := strings.Split(sections[len(sections)-1], ".")
	return nameParts[0]
}

func (ct *CombinedTrace) PrintCombinedInfo(outputDir string) {
	tracefile, err := os.Open(ct.Filepath)
	defer tracefile.Close()
	if err != nil {
		panic(err)
	}
	outfilePath := fmt.Sprintf("%s/%s.csv", outputDir, ct.Filename())
	outfile, err := os.Create(outfilePath)
	defer outfile.Close()

	csvWriter := csv.NewWriter(outfile)
	defer csvWriter.Flush()
	csvWriter.Write(ct.CsvHeaders())
	var state LogLine
	traceScanner := bufio.NewScanner(tracefile)
	firstTime := time.Unix(0, 0)
	for traceScanner.Scan() {
		fmt.Println("got to scan")
		lineBytes := traceScanner.Bytes()
		if err := json.Unmarshal(lineBytes, &state); err != nil {
			panic(err)
		}
		if firstTime.Unix() == 0 {
			firstTime = state.Time.Time
		}
		csvRow := state.ToCsvRow(firstTime)
		csvWriter.Write(csvRow)
	}
}
