package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/aditiharini/drone-monitor/scripts/utils"
	"github.com/knq/hilink"
)

func print(s string) {
	fmt.Println(s)
}

func setupDevices() {
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 down").Run(); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 hw ether 0c:5b:8f:27:9a:65").Run(); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if err := exec.Command("bash", "-c", "sudo ifconfig eth2 up").Run(); err != nil {
		panic(err)
	}
	// time.Sleep(1 * time.Second)
	// if err := exec.Command("bash", "-c", "sudo ifconfig wlan0 down").Run(); err != nil {
	// 	panic(err)
	// }
}

func ifconfig(outfile string) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("ifconfig >> %s", outfile)).CombinedOutput()
	fmt.Println("[ifconfig]", string(out))
	if err != nil {
		panic(err)
	}
}

func main() {
	useTcp := flag.Bool("tcp", false, "iperf protocol")
	flag.Parse()

	startTime := time.Now().Unix()
	ifconfigOutfile := fmt.Sprintf("drone-%d.ifconfig", startTime)
	ifconfig(ifconfigOutfile)

	// Turn saturatr on and off
	count := 0

	proto := "-u -b 30M"
	if *useTcp {
		proto = ""
	}

	outfile := fmt.Sprintf("drone-%d.pcap", startTime)
	tcpdumpCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo tcpdump -n -i eth1 -w %s dst port 5201", outfile))
	utils.RunCmd(tcpdumpCmd, "[tcpdump]", print, print)
	time.Sleep(2 * time.Second)

	time.Sleep(1 * time.Minute)

	client, err := hilink.NewClient(hilink.URL("http://192.168.0.1"))
	if err != nil {
		panic(err)
	}

	logfile, err := os.Create(fmt.Sprintf("drone-%d.hilink", startTime))
	if err != nil {
		panic(err)
	}
	logWriter := bufio.NewWriter(logfile)
	defer logWriter.Flush()

	httpClient := http.Client{Timeout: 1 * time.Second}
	go func() {
		for {
			fmt.Fprintln(logWriter, time.Now().UnixNano())
			info, err := client.TrafficInfo()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintln(logWriter, info)
			info, err = client.NetworkInfo()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintln(logWriter, info)

			info, err = client.SignalInfo()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintln(logWriter, info)
			body, err := json.Marshal(info)
			if err != nil {
				fmt.Println(err)
			}
			res, err := httpClient.Post("http://3.91.1.79:10000/drone/signal", "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err)
				fmt.Println(res)
			}

			info, err = client.StatusInfo()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintln(logWriter, info)

			info, err = client.ModeNetworkInfo()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintln(logWriter, info)
			logWriter.Flush()

			time.Sleep(1 * time.Second)
		}
	}()
	pingUpOutfile := fmt.Sprintf("drone-%d.ping", startTime)
	pingCmd := exec.Command("bash", "-c", fmt.Sprintf("stdbuf -oL ping -i 1 3.91.1.79 | tee %s", pingUpOutfile))
	utils.RunCmd(pingCmd, "[ping]", func(s string) {
		utils.PostLatency(s, httpClient, "http://3.91.1.79:10000/drone/ping")
	}, print)

	for {
		fmt.Printf("[iperf] starting")

		// Do upload
		iperfUploadOutfile := fmt.Sprintf("drone-%d-%d-up.iperf", startTime, count)
		iperfUploadCmd := exec.Command("bash", "-c", fmt.Sprintf("iperf3 %s -t %d -c 3.91.1.79 > %s", proto, 180, iperfUploadOutfile))
		utils.RunCmd(iperfUploadCmd, "[iperf]", print, print)

		iperfUploadCmd.Wait()
		ifconfig(ifconfigOutfile)

		// Do download
		iperfDownloadOutfile := fmt.Sprintf("drone-%d-%d-down.iperf", startTime, count)
		iperfDownloadCmd := exec.Command("bash", "-c", fmt.Sprintf("stdbuf -oL iperf3 -R %s -t %d -c 3.91.1.79 | tee %s", proto, 180, iperfDownloadOutfile))
		utils.RunCmd(iperfDownloadCmd, "[iperf]", func(s string) {
			utils.PostBandwidth(s, "download", httpClient, "http://3.91.1.79:10000/drone/iperf")
		}, print)

		iperfDownloadCmd.Wait()
		ifconfig(ifconfigOutfile)

		count++

		time.Sleep(1 * time.Minute)
		fmt.Printf("[iperf] stopping")
	}
}
