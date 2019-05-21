package kubetool

import (
	"encoding/json"
	"fmt"
	"github.com/wenlaizhou/kubetype"
	"regexp"
)

// 获取节点列表
//
// cluserName 为空, 则获取全部集群node列表
func GetNodes(clusterName string) map[string]kubetype.NodeList {
	res := make(map[string]kubetype.NodeList)
	if len(clusterName) > 0 {
		cmdRes, err := KubeApi(Cluster[clusterName], CmdGet, "no", "-o", "json")
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
		cmdRes, err := KubeApi(c, CmdGet, "no", "-o", "json")
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
	return KubeApi(cluster, CmdGet, "no", node, "-o", "yaml")
}

// 获取节点描述
func DescNode(cluster KubeCluster, node string) (string, error) {
	return KubeApi(cluster, CmdDesc, "no", node, ArgRecursive)
}

// 驱逐节点
func DrainNode(cluster KubeCluster, node string) (string, error) {
	return KubeApi(cluster, CmdDrain, node, "--force=true", "--ignore-daemonsets=true", "--delete-local-data=false")
}

// 删除节点
func DeleteNode(cluster KubeCluster, node string) (string, error) {
	return KubeApi(cluster, CmdDelete, "node", node, "--now")
}

const TraintNodeSpec = "mcloud=deprecated"

const TraintNodeKey = "mcloud"

// 节点增加污点
func TraintNode(cluster KubeCluster, node string) (string, error) {

	// # Update node 'foo' with a taint with key 'dedicated' and value 'special-user' and effect 'NoSchedule'.
	// # If a taint with that key and effect already exists, its value is replaced as specified.
	// kubectl taint nodes foo dedicated=special-user:NoSchedule
	return KubeApi(cluster, Cmdtrait, "nodes", node, fmt.Sprintf("%s:NoSchedule", TraintNodeSpec))

}

// 节点删除污点
func UntraintNode(cluster KubeCluster, node string) (string, error) {
	// # Remove from node 'foo' the taint with key 'dedicated' and effect 'NoSchedule' if one exists.
	// kubectl taint nodes foo dedicated:NoSchedule-
	//
	// # Remove from node 'foo' all the taints with key 'dedicated'
	// kubectl taint nodes foo dedicated-
	return KubeApi(cluster, Cmdtrait, "nodes", node, fmt.Sprintf("%s-", TraintNodeKey))

}

// 节点不可调度
// func UnscheduleNode(cluster KubeCluster, node string) (string, error) {
//
// }

// 节点可调度
// func ScheduleNode(cluster KubeCluster, node string) (string, error) {
//
// }

type NodeResource struct {
	Name    string
	Limit   map[string]string
	Request map[string]string
}

var nodesResourceCache map[KubeCluster]interface{}

var nameReg = regexp.MustCompile("Name:\\s+(\\S+).*")
var cpuResourceReg = regexp.MustCompile("cpu\\s+(\\w+)\\s+\\((\\w+)%?\\)\\s+(\\w+)\\s+\\((\\w+)%?\\).*")
var memResourceReg = regexp.MustCompile("memory\\s+(\\w+)\\s+\\((\\w+)%?\\)\\s+(\\w+)\\s+\\((\\w+)%?\\).*")

// 描述所有节点
func DescAllNodes(cluster KubeCluster) (string, error) {
	return KubeApi(cluster, CmdDesc, "no")
}

// 获取所有节点的资源信息
// todo
func GetAllNodesResource(cluster KubeCluster) []NodeResource {
	return nil
}
