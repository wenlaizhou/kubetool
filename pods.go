package kubetool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"strings"

	"github.com/wenlaizhou/kubetype"
)

// pod轻量列表
type PodResource struct {
	Name     string
	Status   string
	Age      string
	IP       string
	Node     string
	Ready    string
	Restarts string
}

// 获取pod轻量列表
func GetPodsLight(cluster KubeCluster, ns string) []PodResource {
	if len(ns) <= 0 {
		return nil
	}
	args := []string{
		CmdGet, "po",
	}
	args = append(args, "-n", ns, "-o", "wide")
	cmdRes, err := KubeApi(cluster, args...)
	if err != nil {
		return nil
	}
	var result []PodResource
	table := middleware.RenderTable(cmdRes, 7)
	if len(table) <= 0 {
		return result
	}
	for _, row := range table {
		rowData := PodResource{
			Name:     row["NAME"],
			Status:   row["STATUS"],
			Age:      row["AGE"],
			IP:       row["IP"],
			Node:     row["NODE"],
			Ready:    row["READY"],
			Restarts: row["RESTARTS"],
		}
		result = append(result, rowData)
	}
	return result
}

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

// 根据pod属性进行数据查询
// fieldSelector : k=v,k=v
func QueryPods(clusterName string, ns string, fieldSelector string) kubetype.PodList {
	cluster := Cluster[clusterName]
	var args []string
	args = append(args, CmdGet)
	args = append(args, "po")
	args = append(args, "-o")
	args = append(args, "json")
	if len(ns) > 0 {
		args = append(args, "-n", ns)
	} else {
		args = append(args, ArgsAllNamespaces)
	}
	if len(fieldSelector) > 0 {
		args = append(args, fmt.Sprintf(fieldSelectorTpl, fieldSelector))
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

// pod 查询
func QueryPodsByLabel(clusterName string, ns string, selector string) kubetype.PodList {
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

func GetPodByIp(clusterName string, ip string) (kubetype.Pod, error) {
	var res kubetype.Pod
	cluster, hasCluster := Cluster[clusterName]
	if !hasCluster {
		return res, errors.New("cluster is unavailable")
	}
	var args []string
	args = append(args, CmdGet)
	args = append(args, "po")
	args = append(args, ArgsAllNamespaces)
	args = append(args, "-o")
	args = append(args, "json")
	cmdRes, err := KubeApi(cluster, args...)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get pods error : %s", err.Error())
		return res, err
	}
	podList := kubetype.PodList{}
	err = json.Unmarshal([]byte(cmdRes), &podList)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get pods error : %s", err.Error())
		return res, err
	}
	if len(podList.Items) <= 0 {
		return res, errors.New("cluster is unavailable")
	}
	for _, po := range podList.Items {
		if po.Status.PodIP == ip {
			return po, nil
		}
	}
	return res, errors.New("not found")
}

// 获取pod配置信息
func GetPod(cluster KubeCluster, pod string, ns string) (string, error) {
	return KubeApi(cluster, "get", "po", pod, "-n", ns, "-o", "yaml")
}

// 获取pod详细信息
func GetPodItem(cluster KubeCluster, pod string, ns string) (kubetype.Pod, error) {
	res := kubetype.Pod{}
	cmdRes, err := KubeApi(cluster, "get", "po", pod, "-n", ns, "-o", "json")
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(cmdRes), &res)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get pods error : %s", err.Error())
		return res, err
	}
	return res, nil
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

// 获取该节点下的所有pod
func GetPodsByNode(cluster KubeCluster, node string) (kubetype.PodList, error) {
	// 	kubectl get po --field-selector='spec.nodeName=idc02-sre-kubernetes-04' --all-namespaces -o json
	var result kubetype.PodList
	if len(node) <= 0 {
		return result, errors.New("node 为空")
	}
	args := []string{CmdGet, "po", ArgsAllNamespaces, "-o", "json"}
	args = append(args, fmt.Sprintf("--field-selector=spec.nodeName=%s", node))
	res, err := KubeApi(cluster, args...)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(res), &result)
	return result, err
}
