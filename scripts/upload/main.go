package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	uploadCmd := fmt.Sprintf("/home/ubuntu/dropbox_uploader.sh upload %s %s", from, to)
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

func main() {
	trace := flag.String("trace", "", "trace to be processed")
	name := flag.String("name", "", "folder to upload to")
	flag.Parse()
	allIds := allSenderIds(*trace)
	uploadDir := fmt.Sprintf("Drone-Project/measurements/saturatr_traces/%s", *name)
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
	for _, id := range allIds {
		out, err := exec.Command("python", "prep-for-simulation.py", *trace, id).CombinedOutput()
		if err != nil {
			print(string(out))
			panic(err)
		}
		move(fmt.Sprintf("uplink-%s.pps", id), "tmp/processed/traces")
		latencyFile := fmt.Sprintf("tmp/processed/stats/latency-%s.csv", id)
		throughputFile := fmt.Sprintf("tmp/processed/stats/throughput-%s.csv", id)
		traces.WriteCsvData(*trace, id, latencyFile, throughputFile)
	}

	copy(*trace, "tmp/raw")
	if *name != "" {
		upload("tmp", uploadDir)
	}
}
