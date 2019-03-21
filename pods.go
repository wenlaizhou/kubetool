package kubetool

import (
	"encoding/json"

	"github.com/wenlaizhou/kubetype"
)

// 获取pod列表
//
// clusterName为空, 则获取全部集群pod列表
func GetPods(clusterName string) map[string]kubetype.PodList {
	res := make(map[string]kubetype.PodList)
	if len(clusterName) > 0 {
		cmdRes, err := ExecKubectl(Cluster[clusterName], CmdGet, "po", "-o", "json", "--all-namespaces")
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get pods error : %s", clusterName, err.Error())
			return res
		}
		clusterNodes := kubetype.PodList{}
		err = json.Unmarshal([]byte(cmdRes), &clusterNodes)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get pods error : %s", clusterName, err.Error())
			return res
		}
		res[clusterName] = clusterNodes
		return res
	}
	for n, c := range Cluster {
		cmdRes, err := ExecKubectl(c, CmdGet, "po", "-o", "json", "--all-namespaces")
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get pods error : %s", n, err.Error())
			continue
		}
		clusterNodes := kubetype.PodList{}
		err = json.Unmarshal([]byte(cmdRes), &clusterNodes)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get pods error : %s", n, err.Error())
			continue
		}
		res[n] = clusterNodes
	}
	return res
}

// 获取pod配置信息
func GetPod(cluster KubeCluster, pod string, ns string) (string, error) {
	return ExecKubectl(cluster, "get", "po", pod, "-n", ns, "-o", "yaml")
}

// 描述pod
func DescPod(cluster KubeCluster, pod string, ns string) (string, error) {
	return ExecKubectl(cluster, "describe", "po", pod, "-n", ns, "--recursive=true")
}

// 执行pod内部命令
func ExecPodContainer(cluster KubeCluster, pod string, ns string, containerName string, command []string) (string, error) {
	// k exec openresty-proxy-5c5c498949-4mcn5 -n b1 -c openresty-proxy -i -- echo hello world
	args := []string{CmdExec}
	args = append(args, pod)
	if len(ns) > 0 {
		args = append(args, "-n", ns)
	}
	if len(containerName) > 0 {
		args = append(args, "-c", containerName)
	}
	args = append(args, "--")
	args = append(args, command...)
	return ExecKubectl(cluster, args...)
}
