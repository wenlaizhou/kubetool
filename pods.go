package kubetool

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wenlaizhou/kubetype"
)

// 获取pod列表
//
// clusterName为空, 则获取全部集群pod列表
func GetPods(clusterName string, ns string) map[string]kubetype.PodList {
	res := make(map[string]kubetype.PodList)
	args := []string{
		CmdGet, "po", "-o", "json",
	}
	if len(ns) > 0 {
		args = append(args, "-n", ns)
	} else {
		args = append(args, "--all-namespaces")
	}
	if len(clusterName) > 0 {
		cmdRes, err := KubeApi(Cluster[clusterName], args...)
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
		cmdRes, err := KubeApi(c, args...)
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

// pod 查询
func QueryPods(clusterName string, ns string, selector string) kubetype.PodList {
	cluster := Cluster[clusterName]
	var args []string
	args = append(args, CmdGet)
	args = append(args, "po")
	args = append(args, "-o")
	args = append(args, "json")
	args = append(args, "-n")
	if len(ns) <= 0 {
		ns = "default"
	}
	args = append(args, ns)
	if len(selector) > 0 {
		args = append(args, fmt.Sprintf(selectorTpl, selector))
	}
	res := kubetype.PodList{}
	cmdRes, err := KubeApi(cluster, args...)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get pods error : %s", err.Error())
		return res
	}
	err = json.Unmarshal([]byte(cmdRes), &res)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get pods error : %s", err.Error())
	}
	return res
}

// 获取pod配置信息
func GetPod(cluster KubeCluster, pod string, ns string) (string, error) {
	return KubeApi(cluster, "get", "po", pod, "-n", ns, "-o", "yaml")
}

// 描述pod
func DescPod(cluster KubeCluster, pod string, ns string) (string, error) {
	return KubeApi(cluster, "describe", "po", pod, "-n", ns, "--recursive=true")
}

// 执行pod内部命令
func ExecPodContainer(cluster KubeCluster, pod string, ns string, containerName string, workDir string, command []string) (string, error) {
	// k exec openresty-proxy-5c5c498949-4mcn5 -n b1 -c openresty-proxy -i -- echo hello world
	// kubectl exec dubbo-admin-d97dbfbbd-jrjb2 -n b1 -- sh -c "cd /usr && pwd && ping 10.2.3.4"
	args := []string{CmdExec}
	args = append(args, pod)
	if len(ns) > 0 {
		args = append(args, "-n", ns)
	}
	if len(containerName) > 0 {
		args = append(args, "-c", containerName)
	}
	args = append(args, "--")
	if len(workDir) > 0 {
		args = append(args, "sh")
		args = append(args, "-c")
		args = append(args, fmt.Sprintf("cd %s && %s", workDir, strings.Join(command, " ")))
	} else {
		args = append(args, command...)
	}
	return KubeApi(cluster, args...)
}
