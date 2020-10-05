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
	Filepath  string
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
	Time    string
	Rsrp    string
	Rsrq    string
	Rssi    string
	Sinr    string
	Cell_id string
}

func (ht *HilinkTrace) CsvHeader() []string {
	return []string{"time", "cell_id", "rsrp", "rsrq", "rssi", "sinr"}
}

func (s *Signal) ToCsvRow() []string {
	return []string{s.Time, s.CellId, s.Rsrp, s.Rsrq, s.Rssi, s.Sinr}
}

func (rs *RawSignal) ToProcessed() Signal {
	var signal Signal
	signal.CellId = rs.Cell_id
	signal.Rsrp = rs.Rsrp[:len(rs.Rsrp)-3]
	signal.Rsrq = rs.Rsrq[:len(rs.Rsrq)-2]
	signal.Rssi = rs.Rssi[:len(rs.Rssi)-3]
	signal.Sinr = rs.Sinr[:len(rs.Sinr)-2]
	signal.Time = rs.Time
	return signal
}

func (ht *HilinkTrace) ParseChunk(reader *bufio.Scanner) LogChunk {
	timeLine := reader.Text()
	var rawSignal RawSignal
	reader.Scan()
	reader.Scan()
	reader.Scan()
	line := reader.Text()
	line = line[4 : len(line)-1]
	fmt.Println(line)
	pairs := strings.Split(line, " ")
	for _, pair := range pairs {
		splitPair := strings.Split(pair, ":")
		key, value := splitPair[0], splitPair[1]
		field := reflect.ValueOf(&rawSignal).Elem().FieldByName(strings.Title(key))
		if field.CanSet() {
			field.SetString(value)
		}
	}
	reader.Scan()
	reader.Scan()
	rawSignal.Time = timeLine
	var logChunk LogChunk
	logChunk.Signal = rawSignal.ToProcessed()
	return logChunk
}

func (ht *HilinkTrace) Filename() string {
	sections := strings.Split(ht.Filepath, "/")
	nameParts := strings.Split(sections[len(sections)-1], ".")
	return nameParts[0]
}

func (ht *HilinkTrace) PrintSignalInfo(outputDir string) {
	tracefile, err := os.Open(ht.Filepath)
	defer tracefile.Close()
	if err != nil {
		panic(err)
	}
	outfilePath := fmt.Sprintf("%s/%s.csv", outputDir, ht.Filename())
	fmt.Println(outfilePath)
	outfile, err := os.Create(outfilePath)
	defer outfile.Close()
	csvWriter := csv.NewWriter(outfile)
	defer csvWriter.Flush()
	csvWriter.Write(ht.CsvHeader())
	traceReader := bufio.NewScanner(tracefile)
	for traceReader.Scan() {
		chunk := ht.ParseChunk(traceReader)
		csvWriter.Write(chunk.Signal.ToCsvRow())
	}
}
