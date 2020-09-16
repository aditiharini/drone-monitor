package main

import (
	"bufio"
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
	setupDevices()
	time.Sleep(10 * time.Second)
	out, err := exec.Command("ifconfig").CombinedOutput()
	fmt.Println("[ifconfig]", string(out))
	if err != nil {
		panic(err)
	}

	interestmapsDir := "/home/pi/interestmaps/"
	interestmapsCmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s; go run launch.go interestmaps-real.cfg 0 drone", interestmapsDir))
	run(interestmapsCmd, "[interestmaps]", true, true)
	time.Sleep(20 * time.Second)

	// Turn saturatr on and off
	for {
		fmt.Printf("[iperf] starting")
		iperfCmd := exec.Command("bash", "-c", fmt.Sprintf("iperf3 -u -b 25M -t %d -c 3.91.1.79", 60))
		run(iperfCmd, "iperf", true, true)
		iperfCmd.Wait()
		time.Sleep(1 * time.Minute)
		fmt.Printf("[iperf] stopping")
	}
}
