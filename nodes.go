package kubetool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/kubetype"
	"regexp"
	"strings"
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

// 节点描述对象
type NodeDesc struct {
	Name             string
	Status           string
	Roles            string
	Age              string
	Version          string
	InternalIp       string
	ExternalIp       string
	Os               string
	Kernel           string
	ContainerRuntime string
}

var nodesResourceCache map[KubeCluster]interface{}

var nameReg = regexp.MustCompile("Name:\\s+(\\S+).*")
var cpuResourceReg = regexp.MustCompile("cpu\\s+(\\w+)\\s+\\((\\w+)%?\\)\\s+(\\w+)\\s+\\((\\w+)%?\\).*")
var memResourceReg = regexp.MustCompile("memory\\s+(\\w+)\\s+\\((\\w+)%?\\)\\s+(\\w+)\\s+\\((\\w+)%?\\).*")

func DescNodes(clusterName string) ([]NodeDesc, error) {
	cluster, success := Cluster[clusterName]
	if !success {
		return nil, errors.New("不存在该集群")
	}
	res, err := KubeApi(cluster, "get", "no", "-o", "wide")
	if err != nil {
		return nil, err
	}
	var result []NodeDesc
	res = strings.TrimSpace(res)
	nodes := strings.Split(res, "\n")
	nodes = nodes[1:]
	for _, node := range nodes {
		fields := strings.Fields(node)
		nodeDesc := NodeDesc{}
		// NAME STATUS ROLES AGE VERSION INTERNAL-IP EXTERNAL-IP OS-IMAGE KERNEL-VERSION CONTAINER-RUNTIME
		nodeDesc.Name = fields[0]
		nodeDesc.Status = fields[1]
		nodeDesc.Roles = fields[2]
		nodeDesc.Age = fields[3]
		nodeDesc.Version = fields[4]
		nodeDesc.InternalIp = fields[5]
		nodeDesc.ExternalIp = fields[6]
		nodeDesc.Os = fields[7]
		nodeDesc.Kernel = fields[8]
		nodeDesc.ContainerRuntime = fields[9]
		result = append(result, nodeDesc)
	}
	return result, nil
}

// 获取集群全部节点描述
func DescAllNodes(clusterName string) (kubetype.NodeList, error) {
	result := kubetype.NodeList{}
	cluster, success := Cluster[clusterName]
	if !success {
		return result, errors.New("不存在该集群")
	}
	res, err := KubeApi(cluster, "get", "no", "-o", "wide")
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(strings.TrimSpace(res)), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// // 描述所有节点
// func DescAllNodes() (map[string]string, error) {
// 	result := make(map[string]string)
// 	for _, cluster := range Cluster {
// 		desc, err := DescClusterNodes(cluster)
// 		if err != nil {
// 			return result, err
// 		}
// 		result[cluster.Name] = desc
// 	}
// 	return result, nil
// }

// 描述集群内的全部节点
func DescClusterNodes(cluster KubeCluster) (string, error) {
	return KubeApi(cluster, CmdDesc, "no")
}

// 获取所有节点的资源信息
// todo
func GetAllNodesResource(cluster KubeCluster) []NodeResource {
	return nil
}
