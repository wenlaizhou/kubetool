package kubetool

// 获取pod日志
func GetLog(cluster KubeCluster, pod string, ns string) (string, error) {
	return KubeApi(cluster, CmdLogs, ArgsAllContainers, pod, ArgsNamespace, ns, "--tail=200")
}

// 获取pod所有日志
func GetAllLog(cluster KubeCluster, pod string, ns string) (string, error) {
	return KubeApi(cluster, CmdLogs, ArgsAllContainers, pod, ArgsNamespace, ns)
}
