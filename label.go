package kubetool

import "fmt"

// 创建标签
func AddLabel(cluster KubeCluster, resourceName string, name string, key string, value string, ns string) (string, error) {
	var args []string
	args = append(args, CmdLabel)
	args = append(args, ArgsOverwrite)
	args = append(args, resourceName)
	args = append(args, name)
	args = append(args, fmt.Sprintf("%s=%s", key, value))
	if len(ns) > 0 {
		args = append(args, "-n")
		args = append(args, ns)
	}
	return ExecKubectl(cluster, args...)
}

// 删除标签
func DeleteLabel(cluster KubeCluster, resourceName string, name string, key string, ns string) (string, error) {
	var args []string
	args = append(args, CmdLabel)
	args = append(args, ArgsOverwrite)
	args = append(args, resourceName)
	args = append(args, name)
	args = append(args, fmt.Sprintf("%s-", key))
	if len(ns) > 0 {
		args = append(args, "-n")
		args = append(args, ns)
	}
	return ExecKubectl(cluster, args...)
}
