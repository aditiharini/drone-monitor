package trace

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aditiharini/drone-monitor/api"
)

type CombinedTrace struct {
	OutputDir string
	Filepath  string
}

type LogLine struct {
	Time  time.Time `json:"time"`
	State api.State `json:"state"`
}

func (ll *LogLine) ToCsvRow() []string {
	timeStr := strconv.FormatInt(ll.Time.Unix(), 10)
	rsrp := ll.State.Drone.Signal.Rsrp
	if rsrp == "-1" {
		rsrp = "NA"
	}
	rsrq := ll.State.Drone.Signal.Rsrq
	if rsrq == "-1" {
		rsrq = "NA"
	}
	rssi := ll.State.Drone.Signal.Rssi
	if rssi == "-1" {
		rssi = "NA"
	}
	sinr := ll.State.Drone.Signal.Sinr
	if sinr == "-1" {
		sinr = "NA"
	}
	uplinkBandwidth := fmt.Sprintf("%f", ll.State.Server.Iperf.Bandwidth)
	if ll.State.Server.Iperf.Bandwidth == -1 {
		uplinkBandwidth = "NA"
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
	return []string{timeStr, rsrp, rsrq, rssi, sinr, uplinkBandwidth, downlinkBandwidth, latency}
}

func (ct *CombinedTrace) CsvHeadders() []string {
	return []string{"time", "rsrp", "rsrq", "rssi", "sinr", "uplinkBw", "downlinkBw", "latency"}
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
	var state LogLine
	traceScanner := bufio.NewScanner(tracefile)
	for traceScanner.Scan() {
		lineBytes := traceScanner.Bytes()
		json.Unmarshal(lineBytes, state)
		csvRow := state.ToCsvRow()
		csvWriter.Write(csvRow)
	}
}
