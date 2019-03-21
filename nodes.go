package kubetool

import (
	"encoding/json"
	"fmt"
	"github.com/wenlaizhou/kubetype"
)

// 获取节点列表
//
// cluserName 为空, 则获取全部集群node列表
func GetNodes(clusterName string) map[string]kubetype.NodeList {
	res := make(map[string]kubetype.NodeList)
	if len(clusterName) > 0 {
		cmdRes, err := ExecKubectl(Cluster[clusterName], CmdGet, "no", "-o", "json")
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get node error : %s", clusterName, err.Error())
			return res
		}
		clusterNodes := kubetype.NodeList{}
		err = json.Unmarshal([]byte(cmdRes), &clusterNodes)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get node error : %s", clusterName, err.Error())
			return res
		}
		res[clusterName] = clusterNodes
		return res
	}
	for n, c := range Cluster {
		cmdRes, err := ExecKubectl(c, CmdGet, "no", "-o", "json")
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get node error : %s", n, err.Error())
			continue
		}
		clusterNodes := kubetype.NodeList{}
		err = json.Unmarshal([]byte(cmdRes), &clusterNodes)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get node error : %s", n, err.Error())
			continue
		}
		res[n] = clusterNodes
	}
	return res
}

// 获取节点
func GetNode(cluster KubeCluster, node string) (string, error) {
	return ExecKubectl(cluster, CmdGet, "no", node, "-o", "yaml")
}

// 获取节点描述
func DescNode(cluster KubeCluster, node string) (string, error) {
	return ExecKubectl(cluster, CmdDesc, "no", node, ArgRecursive)
}

// 驱逐节点
func DrainNode(cluster KubeCluster, node string) (string, error) {
	return ExecKubectl(cluster, CmdDrain, node, "--force=true", "--ignore-daemonsets=true", "--delete-local-data=false")
}

// 删除节点
func DeleteNode(cluster KubeCluster, node string) (string, error) {
	return ExecKubectl(cluster, CmdDelete, "node", node, "--now")
}

const TraintNodeSpec = "mcloud=deprecated"

const TraintNodeKey = "mcloud"

// 节点增加污点
func TraintNode(cluster KubeCluster, node string) (string, error) {

	// # Update node 'foo' with a taint with key 'dedicated' and value 'special-user' and effect 'NoSchedule'.
	// # If a taint with that key and effect already exists, its value is replaced as specified.
	// kubectl taint nodes foo dedicated=special-user:NoSchedule
	return ExecKubectl(cluster, Cmdtrait, "nodes", node, fmt.Sprintf("%s:NoSchedule", TraintNodeSpec))

}

// 节点删除污点
func UntraintNode(cluster KubeCluster, node string) (string, error) {
	// # Remove from node 'foo' the taint with key 'dedicated' and effect 'NoSchedule' if one exists.
	// kubectl taint nodes foo dedicated:NoSchedule-
	//
	// # Remove from node 'foo' all the taints with key 'dedicated'
	// kubectl taint nodes foo dedicated-
	return ExecKubectl(cluster, Cmdtrait, "nodes", node, fmt.Sprintf("%s-", TraintNodeKey))

}
