package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	trace "github.com/aditiharini/drone-monitor/scripts/traces"
	traces "github.com/aditiharini/drone-monitor/scripts/traces"
)

func allSenderIds(filename string) []string {
	ids := make(map[string]bool)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "INCOMING DATA RECEIVED") {
			id := strings.Split(strings.Split(strings.Split(line, ", ")[0], " ")[3], "=")[1]
			ids[id] = true
		}
	}

	var idArr []string
	for id := range ids {
		idArr = append(idArr, id)
	}
	return idArr
}

func upload(from string, to string) {
	uploadCmd := fmt.Sprintf("dropbox_uploader.sh upload %s %s", from, to)
	out, err := exec.Command("bash", "-c", uploadCmd).CombinedOutput()
	if err != nil {
		print(string(out))
		panic(err)
	}
}

func move(from string, to string) {
	moveCmd := fmt.Sprintf("mv %s %s", from, to)
	out, err := exec.Command("bash", "-c", moveCmd).CombinedOutput()
	if err != nil {
		print(string(out))
		panic(err)
	}
}

func copy(from string, to string) {
	copyCmd := fmt.Sprintf("cp %s %s", from, to)
	out, err := exec.Command("bash", "-c", copyCmd).CombinedOutput()
	if err != nil {
		print(string(out))
		panic(err)
	}
}

type DirectoryStructure struct {
	name     string
	children []DirectoryStructure
}

func (ds DirectoryStructure) create() {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("rm -rf %s", ds.name)).CombinedOutput()
	if err != nil {
		print(string(out))
		panic(err)
	}

	out, err = exec.Command("bash", "-c", fmt.Sprintf("mkdir %s", ds.name)).CombinedOutput()
	if err != nil {
		print(string(out))
		panic(err)
	}
	for _, child := range ds.children {
		child.name = fmt.Sprintf("%s/%s", ds.name, child.name)
		child.create()
	}
}

type TraceMetadata struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Weather string `json:"weather"`
	Notes   string `json:"notes"`
}

func main() {
	traceFile := flag.String("trace", "", "trace to be processed")
	name := flag.String("name", "", "folder to upload to")
	traceType := flag.String("type", "iperf", "type of trace")
	isTcp := flag.Bool("tcp", false, "collection used tcp")
	flag.Parse()

	dirStructure := DirectoryStructure{
		name: "tmp",
		children: []DirectoryStructure{
			{name: "raw"},
			{
				name: "processed",
				children: []DirectoryStructure{
					{name: "stats"},
					{name: "traces"},
				},
			},
		},
	}
	dirStructure.create()

	var uploadDir string
	if *traceType == "saturatr" {
		allIds := allSenderIds(*traceFile)
		uploadDir = fmt.Sprintf("Drone-Project/measurements/saturatr_traces/%s", *name)
		for _, id := range allIds {
			out, err := exec.Command("python", "prep-for-simulation.py", *traceFile, id).CombinedOutput()
			if err != nil {
				print(string(out))
				panic(err)
			}
			move(fmt.Sprintf("uplink-%s.pps", id), "tmp/processed/traces")
			latencyFile := fmt.Sprintf("tmp/processed/stats/latency-%s.csv", id)
			throughputFile := fmt.Sprintf("tmp/processed/stats/throughput-%s.csv", id)
			traces.WriteCsvData(*traceFile, id, latencyFile, throughputFile)
		}
	} else if *traceType == "iperf" {
		uploadDir = fmt.Sprintf("Drone-Project/measurements/iperf_traces/%s", *name)
		proto := "udp"
		if *isTcp {
			proto = "tcp"
		}
		processor := trace.PcapProcessor{Filename: *traceFile, OutputDir: "tmp/processed/traces", Filter: proto, CurrentFilenum: 0}
		processor.ToMahiMahi()
		for _, file := range processor.MahimahiFiles {
			mmTrace := trace.MahimahiTrace{Filename: file, Dirname: "tmp/processed/traces", PacketSize: 1500}
			mmTrace.PrintBandwidth("tmp/processed/stats")
		}

	} else if *traceType == "hilink" {
		uploadDir = fmt.Sprintf("Drone-Project/measurements/hilink_traces/%s", *name)
		hilinkTrace := trace.HilinkTrace{Filepath: *traceFile, OutputDir: "tmp/processed/traces"}
		hilinkTrace.PrintSignalInfo("tmp/processed/stats")
	} else {
		panic("invalid trace type")
	}

	copy(*traceFile, "tmp/raw")
	if *name != "" {
		inputScanner := bufio.NewScanner(os.Stdin)
		var metadata TraceMetadata
		fmt.Printf("Date:")
		inputScanner.Scan()
		metadata.Date = inputScanner.Text()
		fmt.Printf("Time:")
		inputScanner.Scan()
		metadata.Time = inputScanner.Text()
		fmt.Printf("Weather:")
		inputScanner.Scan()
		metadata.Weather = inputScanner.Text()
		fmt.Printf("Notes:")
		inputScanner.Scan()
		metadata.Notes = inputScanner.Text()
		fmt.Printf("Finish upload?")
		inputScanner.Scan()
		descFile, err := os.Create("tmp/description.json")
		defer descFile.Close()
		if err != nil {
			panic(err)
		}
		descBytes, err := json.Marshal(metadata)
		if err != nil {
			panic(err)
		}
		_, err = descFile.Write(descBytes)
		if err != nil {
			panic(err)
		}
		if strings.ToLower(inputScanner.Text()) == "y" {
			upload("tmp", uploadDir)
		} else {
			fmt.Println("Aborted")
		}
	}
}
