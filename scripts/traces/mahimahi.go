package trace

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MahimahiTrace struct {
	Dirname    string
	Filename   string
	PacketSize int
}

func (m *MahimahiTrace) PrintBandwidth(outputDir string) {
	stdout := false
	if outputDir == "" {
		stdout = true
	}
	file, err := os.Open(fmt.Sprintf("%s/%s", m.Dirname, m.Filename))
	if err != nil {
		panic(err)
	}
	var tputCsvWriter *csv.Writer
	if !stdout {
		tputCsvName := strings.Split(m.Filename, ".")[0]
		tputCsv, err := os.Create(fmt.Sprintf("%s/%s.csv", outputDir, tputCsvName))
		if err != nil {
			panic(err)
		}
		defer tputCsv.Close()
		tputCsvWriter = csv.NewWriter(tputCsv)
		tputCsvWriter.Write([]string{"time", "throughput"})
		defer tputCsvWriter.Flush()
	}
	scanner := bufio.NewScanner(file)
	lastCalcTime := -1
	numPackets := 0
	for scanner.Scan() {
		curTime, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}

		if lastCalcTime == -1 {
			lastCalcTime = curTime
		}

		if curTime-lastCalcTime >= 1000 {
			numBits := float64(numPackets * m.PacketSize * 8)
			numMbits := numBits / 1000000.
			elapsed := float64(curTime-lastCalcTime) / 1000.
			mbps := numMbits / elapsed
			if stdout {
				fmt.Println(mbps)
			} else {
				tputCsvWriter.Write([]string{strconv.Itoa(lastCalcTime), strconv.Itoa(int(mbps))})
			}
			numPackets = 0
			lastCalcTime = curTime
		}
		numPackets++
	}
}
