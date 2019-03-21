package kubetool

// 获取pod日志
func GetLog(cluster KubeCluster, pod string, ns string) (string, error) {
	return ExecKubectl(cluster, CmdLogs, "--all-containers=true", pod, "-n", ns, "--tail=400")
}
