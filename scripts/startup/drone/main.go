package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/aditiharini/drone-monitor/scripts/utils"
	"github.com/knq/hilink"
)

func run(cmd *exec.Cmd, tag string, printStdout bool, printStderr bool) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if printStdout {
		go func() {
			buf := bufio.NewReader(stdout)
			for {
				line, _, err := buf.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(tag, string(line))
			}
		}()
	}

	if printStderr {
		go func() {
			buf := bufio.NewReader(stderr)
			for {
				line, _, err := buf.ReadLine()
				if err != nil {
					break
				}
				fmt.Println(tag, string(line))
			}
		}()
	}
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

func ifconfig() {
	out, err := exec.Command("ifconfig").CombinedOutput()
	fmt.Println("[ifconfig]", string(out))
	if err != nil {
		panic(err)
	}
}

func main() {
	useTcp := flag.Bool("tcp", false, "iperf protocol")
	flag.Parse()

	ifconfig()

	// Turn saturatr on and off
	startTime := time.Now().Unix()
	count := 0

	proto := "-u -b 30M"
	if *useTcp {
		proto = ""
	}

	outfile := fmt.Sprintf("drone-%d.pcap", startTime)
	tcpdumpCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo tcpdump -n -i eth1 -w %s dst port 5201", outfile))
	utils.RunCmd(tcpdumpCmd, "[tcpdump]", true, true)
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

	go func() {
		for {
			info, err := client.TrafficInfo()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(logWriter, info)
			info, err = client.NetworkInfo()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(logWriter, info)

			info, err = client.SignalInfo()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(logWriter, info)

			info, err = client.StatusInfo()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(logWriter, info)

			info, err = client.ModeNetworkInfo()
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(logWriter, info)
			logWriter.Flush()

			time.Sleep(2 * time.Second)
		}
	}()

	for {
		fmt.Printf("[iperf] starting")

		// Do upload
		iperfUploadOutfile := fmt.Sprintf("%d-%d-up.iperf", startTime, count)
		iperfUploadCmd := exec.Command("bash", "-c", fmt.Sprintf("iperf3 %s -t %d -c 3.91.1.79 > %s", proto, 180, iperfUploadOutfile))
		run(iperfUploadCmd, "iperf", true, true)

		pingUpOutfile := fmt.Sprintf("%d-%d-up.ping", startTime, count)
		pingCmd := exec.Command("bash", "-c", fmt.Sprintf("ping -i 1 3.91.1.79 > %s", pingUpOutfile))
		run(pingCmd, "ping", true, true)

		iperfUploadCmd.Wait()
		pingCmd.Process.Kill()
		ifconfig()

		// Do download
		iperfDownloadOutfile := fmt.Sprintf("%d-%d-down.iperf", startTime, count)
		iperfDownloadCmd := exec.Command("bash", "-c", fmt.Sprintf("iperf3 -R %s -t %d -c 3.91.1.79 > %s", proto, 180, iperfDownloadOutfile))
		run(iperfDownloadCmd, "iperf", true, true)

		pingDownOutfile := fmt.Sprintf("%d-%d-down.ping", startTime, count)
		pingCmd = exec.Command("bash", "-c", fmt.Sprintf("ping -i 1 3.91.1.79 > %s", pingDownOutfile))
		run(pingCmd, "ping", true, true)

		iperfDownloadCmd.Wait()
		pingCmd.Process.Kill()
		ifconfig()

		count++

		time.Sleep(1 * time.Minute)
		fmt.Printf("[iperf] stopping")
	}
}
