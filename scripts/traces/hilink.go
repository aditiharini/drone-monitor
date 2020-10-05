package trace

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type HilinkTrace struct {
	OutputDir string
	Filename  string
}

type LogChunk struct {
	Signal Signal
}

type Signal struct {
	Time   string
	Rsrp   string
	Rsrq   string
	Rssi   string
	Sinr   string
	CellId string
}

type RawSignal struct {
	time    string
	rsrp    string
	rsrq    string
	rssi    string
	sinr    string
	cell_id string
}

func (ht *HilinkTrace) CsvHeader() []string {
	return []string{"time", "cell_id", "rsrp", "rsrq", "rssi", "sinr"}
}

func (s *Signal) ToCsvRow() []string {
	return []string{s.Time, s.CellId, s.Rsrp, s.Rsrq, s.Rssi, s.Sinr}
}

func (rs *RawSignal) ToProcessed() Signal {
	var signal Signal
	signal.CellId = rs.cell_id
	signal.Rsrp = rs.rsrp[:len(rs.rsrp)-3]
	signal.Rsrq = rs.rsrq[:len(rs.rsrq)-2]
	signal.Rssi = rs.rssi[:len(rs.rssi)-3]
	signal.Sinr = rs.sinr[:len(rs.sinr)-2]
	signal.Time = rs.time
	return signal
}

func (ht *HilinkTrace) ParseChunk(reader *bufio.Scanner) LogChunk {
	timeLine := reader.Text()
	var rawSignal RawSignal
	for i := 0; i < 5; i++ {
		line := reader.Text()
		line = line[4 : len(line)-1]
		pairs := strings.Split(line, " ")
		for _, pair := range pairs {
			splitPair := strings.Split(pair, ":")
			key, value := splitPair[0], splitPair[1]
			s := reflect.ValueOf(&rawSignal).Elem()
			s.FieldByName(key).SetString(value)
		}
	}
	rawSignal.time = timeLine
	var logChunk LogChunk
	logChunk.Signal = rawSignal.ToProcessed()
	return logChunk
}

func (ht *HilinkTrace) PrintSignalInfo(outputDir string) {
	tracefile, err := os.Open(ht.Filename)
	if err != nil {
		panic(err)
	}

	outfilePath := fmt.Sprintf("%s/%s", outputDir, ht.Filename)
	outfile, err := os.Create(outfilePath)
	csvWriter := csv.NewWriter(outfile)
	csvWriter.Write(ht.CsvHeader())
	traceReader := bufio.NewScanner(tracefile)
	for traceReader.Scan() {
		chunk := ht.ParseChunk(traceReader)
		csvWriter.Write(chunk.Signal.ToCsvRow())
	}
}
