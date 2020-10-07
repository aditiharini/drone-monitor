package main

import (
	"fmt"
	"os/exec"
	"time"

	trace "github.com/aditiharini/drone-monitor/scripts/traces"
	"github.com/aditiharini/drone-monitor/scripts/utils"
)

func print(s string) {
	fmt.Println(s)
}

func main() {
	monitorDir := "/home/ubuntu/data/drone-monitor/"
	fmt.Println("Starting drone monitor")
	monitorCmd := exec.Command(monitorDir + "drone-monitor")
	utils.RunCmd(monitorCmd, "[monitor]", print, print)
	fmt.Println("Starting interestmaps")
	time.Sleep(2 * time.Second)

	frontendDir := "/home/ubuntu/data/drone-monitor/client/simple/build"
	frontendCmd := exec.Command("serve", frontendDir)
	utils.RunCmd(frontendCmd, "[frontend]", print, print)
	time.Sleep(2 * time.Second)

	interestmapsDir := "/home/ubuntu/interestmaps/"
	interestmapsCmd := exec.Command(interestmapsDir+"interestmaps", interestmapsDir+"interestmaps-real.cfg")
	utils.RunCmd(interestmapsCmd, "[interestmaps]", print, print)
	time.Sleep(2 * time.Second)

	out, err := exec.Command("bash", "-c", "sudo ethtool -K eth0 gro off").CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	iperfOutfile := fmt.Sprintf("%d.iperf", time.Now().Unix())
	iperfCmd := exec.Command("bash", "-c", fmt.Sprintf("stdbuf -oL iperf3 -s | tee %s", iperfOutfile))
	utils.RunCmd(iperfCmd, "[iperf]", func(s string) {
		trace.PostBandwidth(s, "http://3.91.1.79:10000/server/iperf")
	}, print)
	fmt.Println("Starting iperf")
	time.Sleep(2 * time.Second)

	outfile := fmt.Sprintf("srv-%d.pcap", time.Now().Unix())
	tcpdumpCmd := exec.Command("bash", "-c", fmt.Sprintf("sudo tcpdump -n -i eth0 -w %s dst port 5201", outfile))
	utils.RunCmd(tcpdumpCmd, "[tcpdump]", print, print)
	time.Sleep(2 * time.Second)

	iperfCmd.Wait()
	tcpdumpCmd.Wait()
	monitorCmd.Wait()
	frontendCmd.Wait()
	interestmapsCmd.Wait()
}
