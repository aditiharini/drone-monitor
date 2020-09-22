package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/aditiharini/drone-monitor/scripts/utils"
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

func main() {
	monitorDir := "/home/ubuntu/data/drone-monitor/"
	fmt.Println("Starting drone monitor")
	monitorCmd := exec.Command(monitorDir + "drone-monitor")
	run(monitorCmd, "[monitor]", true, true)
	fmt.Println("Starting interestmaps")
	time.Sleep(2 * time.Second)

	frontendDir := "/home/ubuntu/data/drone-monitor/client/simple/build"
	frontendCmd := exec.Command("serve", frontendDir)
	run(frontendCmd, "[frontend]", true, true)
	time.Sleep(2 * time.Second)

	interestmapsDir := "/home/ubuntu/interestmaps/"
	interestmapsCmd := exec.Command(interestmapsDir+"interestmaps", interestmapsDir+"interestmaps-real.cfg")
	run(interestmapsCmd, "[interestmaps]", true, true)
	time.Sleep(2 * time.Second)

	out, err := exec.Command("bash", "-c", "sudo ethtool -K eth0 gro off").CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	iperfOutfile := fmt.Sprintf("%d.iperf", time.Now().Unix())
	iperfCmd := exec.Command("bash", "-c", fmt.Sprintf("stdbuf -oL iperf3 -s | tee %s", iperfOutfile))
	run(iperfCmd, "[iperf]", true, true)
	fmt.Println("Starting iperf")
	time.Sleep(2 * time.Second)

	outfile := fmt.Sprintf("srv-%d.pcap", time.Now().Unix())
	tcpdumpCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo tcpdump -n -i eth0 -w %s dst port 5201", outfile))
	utils.RunCmd(tcpdumpCmd, "[tcpdump]", true, true)
	time.Sleep(2 * time.Second)

	iperfCmd.Wait()
	tcpdumpCmd.Wait()
	monitorCmd.Wait()
	frontendCmd.Wait()
	interestmapsCmd.Wait()
}
