package trace

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type Line struct {
	id          string
	sendTime    int64
	receiveTime int64
	delay       int
}

func parseLine(line string) Line {
	line = strings.TrimSpace(line)
	if strings.Contains(line, "INCOMING DATA RECEIVED") {
		separatedInfo := strings.Split(line, ", ")
		id := strings.Split(strings.Split(separatedInfo[0], " ")[3], "=")[1]
		sendTime, err := strconv.ParseInt(strings.Split(separatedInfo[2], "=")[1], 10, 64)
		if err != nil {
			panic(err)
		}
		receiveTime, err := strconv.ParseInt(strings.Split(separatedInfo[3], "=")[1], 10, 64)
		if err != nil {
			panic(err)
		}
		delay, err := strconv.ParseFloat(strings.Split(separatedInfo[4], "=")[1], 64)
		if err != nil {
			panic(err)
		}
		return Line{id: id, sendTime: sendTime, receiveTime: receiveTime, delay: int(delay * 1000)}
	}
	return Line{id: "-1"}
}

func WriteCsvData(traceFile string, senderId string, latencyFile string, throughputFile string) {
	trace, err := os.Open(traceFile)
	if err != nil {
		panic(err)
	}
	defer trace.Close()

	latency, err := os.Create(latencyFile)
	if err != nil {
		panic(err)
	}
	defer latency.Close()

	throughput, err := os.Create(throughputFile)
	if err != nil {
		panic(err)
	}
	defer throughput.Close()

	latencyWriter := csv.NewWriter(latency)
	throughputWriter := csv.NewWriter(throughput)

	scanner := bufio.NewScanner(trace)
	firstTime := int64(-1)
	prevTime := int64(-1)
	if err := latencyWriter.Write([]string{"time", "latency(ms)"}); err != nil {
		panic(err)
	}
	defer latencyWriter.Flush()

	if err := throughputWriter.Write([]string{"time", "throughput(Mbps)"}); err != nil {
		panic(err)
	}
	defer throughputWriter.Flush()

	numPacketsReceived := 0
	for scanner.Scan() {
		line := parseLine(scanner.Text())
		if line.id != "-1" && line.id == senderId {
			numPacketsReceived++
			if firstTime == -1 {
				firstTime = line.receiveTime
				prevTime = line.receiveTime
			}
			if line.receiveTime-prevTime >= 1000000000 {
				throughputValue := float64(numPacketsReceived*1400.*8./1000000.) / (float64(line.receiveTime-prevTime) / 1000000000.)
				prevTimeOffset := int((prevTime - firstTime) / 1000000)
				latestThroughput := []string{strconv.Itoa(prevTimeOffset), strconv.Itoa(int(throughputValue))}
				if err := throughputWriter.Write(latestThroughput); err != nil {
					panic(err)
				}
				prevTime = line.receiveTime
				numPacketsReceived = 0
			}
			offsetTime := int((line.receiveTime - firstTime) / 1000000)
			data := []string{strconv.Itoa(offsetTime), strconv.Itoa(line.delay)}
			if err := latencyWriter.Write(data); err != nil {
				panic(err)
			}
		}
	}
}
