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
const saveConfig = "--save-config=%s" // true | false A label selector to use for this service.
// Only equality-based selector requirements are supported. If empty (the default) infer the selector
// from the replication controller or replica set.)
const targetPort = "--target-port=%s"
const name = "--name=%s"
const externalIp = "--external-ip=%s"
const selectorTpl = "--selector=%s"

func Expose(cluster KubeCluster, kind string, resourceName string,
	serviceName string, ns string, externalPort string, existPort string,
	nodeIp string,
	selector string) (string, error) {
	// ExternalName - Exposes the Service using an arbitrary name
	// (specified by externalName in the spec) by returning a CNAME
	// record with the name. No proxy is used. This type requires v1.7 or higher of kube-dns.
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
	if len(selector) >= 0 {
		args = append(args, fmt.Sprintf(selectorTpl, selector))
	}
	return KubeApi(cluster, args...)
}
