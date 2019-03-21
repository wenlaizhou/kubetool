package kubetool

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execCmd(cmd string, args ...string) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	fmt.Printf("%s %v\n", cmd, args)
	if err != nil {
		fmt.Printf("%s\n%v", string(output), err)
	} else {
		println(string(output))
	}
}

// 清空节点
func Clean() {

	println("请先执行: kubeadm reset -v 10")
	reader := bufio.NewReader(os.Stdin)

	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	result = strings.TrimSpace(result)
	if result != "yes" && result != "y" {
		println("请输入yes或者y确认执行完毕kubeadm reset")
		return
	}

	// stop service
	println("stop service")
	execCmd("systemctl", "stop", "kubelet")
	execCmd("systemctl", "stop", "docker")

	// delete iptables
	// iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -Xs
	println("delete iptables")
	execCmd("iptables", "-F")
	execCmd("iptables", "-t", "nat", "-F")
	execCmd("iptables", "-t", "mangle", "-F")
	execCmd("iptables", "-X")

	// delete route
	println("delete route")
	output, err := exec.Command("ip", "route").CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "10.244") || strings.Contains(line, "172.17") {
			execCmd("ip", "route", "del", strings.TrimSpace(line))
		}
	}
	// delete iface
	println("delete ifconfig")
	execCmd("ifconfig", "cni0", "down")
	execCmd("ifconfig", "flannel.1", "down")
	execCmd("ifconfig", "docker0", "down")

	// delete bridge
	println("delete bridge")
	execCmd("brctl", "delbr", "cni0")
	execCmd("brctl", "delbr", "docker0")

	// rm files
	println("delete files")
	execCmd("rm", "-rf", "/var/lib/cni")
	execCmd("rm", "-rf", "/var/lib/kubelet")
	execCmd("rm", "-rf", "/etc/cni")
	execCmd("rm", "-rf", "/etc/kubernetes")
}
