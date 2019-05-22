package kubetool

import "fmt"

func Scale(cluster KubeCluster, resourceName string, name string, scaleNumber int, ns string) (string, error) {
	var args []string
	args = append(args, CmdScale)
	args = append(args, fmt.Sprintf("--replicas=%v", scaleNumber))
	args = append(args, fmt.Sprintf("%v/%v", resourceName, name))
	return KubeApi(cluster, args...)
}

const ExposeTypeClusterIP = "ClusterIP" // default
const ExposeTypeNodePort = "NodePort"
const ExposeTypeLoadBalancer = "LoadBalancer"
const ExposeTypeExternalName = "ExternalName"
const typeArg = "--type=%s"
const externalPortTpl = "--port=%s"
const saveConfig = "--save-config=%s" // true | false
const targetPort = "--target-port=%s"
const name = "--name=%s"
const externalIp = "--external-ip=%s"
const selectorTpl = "--selector=%s"

// 对外发布服务
//
// kind: pod (po), service (svc), replicationcontroller (rc), deployment (deploy), replicaset (rs)
//
// resourceName: 当前集群已经存在的资源名称, 可空
//
// serviceName: 对外暴露服务的名称
func Expose(cluster KubeCluster, kind string, resourceName string,
	serviceName string, ns string, externalPort string, existPort string,
	nodeIp string) (string, error) {
	var args []string
	args = append(args, CmdExpose)
	args = append(args, kind)
	args = append(args, resourceName)
	if len(ns) <= 0 {
		ns = "default"
	}
	args = append(args, "-n")
	args = append(args, ns)
	args = append(args, fmt.Sprintf(name, serviceName))
	args = append(args, fmt.Sprintf(targetPort, existPort))
	args = append(args, fmt.Sprintf(externalPortTpl, externalPort))
	args = append(args, fmt.Sprintf(externalIp, nodeIp))
	args = append(args, fmt.Sprintf(typeArg, ExposeTypeNodePort))
	return KubeApi(cluster, args...)
}
