package main

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"time"
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

func main() {
	useTcp := flag.Bool("tcp", false, "iperf protocol")
	flag.Parse()

	out, err := exec.Command("ifconfig").CombinedOutput()
	fmt.Println("[ifconfig]", string(out))
	if err != nil {
		panic(err)
	}

	// Turn saturatr on and off
	startTime := time.Now().Unix()
	count := 0

	proto := "-u"
	if *useTcp {
		proto = ""
	}

	for {
		fmt.Printf("[iperf] starting")
		iperfOutfile := fmt.Sprintf("%d-%d.iperf", startTime, count)
		iperfCmd := exec.Command("bash", "-c", fmt.Sprintf("iperf3 %s -b 30M -t %d -c 3.91.1.79 > %s", proto, 60, iperfOutfile))
		run(iperfCmd, "iperf", true, true)

		pingOutfile := fmt.Sprintf("%d-%d.ping", startTime, count)
		pingCmd := exec.Command("bash", "-c", fmt.Sprintf("ping -i 1 3.91.1.79 > %s", pingOutfile))
		run(pingCmd, "ping", true, true)

		iperfCmd.Wait()
		pingCmd.Process.Kill()
		count++

		time.Sleep(1 * time.Minute)
		fmt.Printf("[iperf] stopping")
	}
}
