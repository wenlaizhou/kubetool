package kubetool

import (
	"fmt"
	"github.com/wenlaizhou/middleware"
	"regexp"
)

func Scale(cluster KubeCluster, resourceName string, name string, scaleNumber int, ns string) (string, error) {
	var args []string
	args = append(args, CmdScale)
	args = append(args, fmt.Sprintf("--replicas=%v", scaleNumber))
	args = append(args, fmt.Sprintf("%v/%v", resourceName, name))
	args = append(args, "-n", ns)
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
const selectorTpl = "-l %s"

// 对外发布服务
//
// kind: pod (po), service (svc), replicationcontroller (rc), deployment (deploy), replicaset (rs)
//
// resourceName: 当前集群已经存在的资源名称, 可空
//
// serviceName: 对外暴露服务的名称
func Expose(cluster KubeCluster, kind string, exposeName string, resourceName string,
	ns string, externalPort string, internalPort string, nodeIp string) (string, error) {

	var args []string
	args = append(args, CmdExpose)
	args = append(args, kind)
	args = append(args, resourceName)
	if len(ns) <= 0 {
		ns = "default"
	}
	serviceName := exposeName
	if len(exposeName) <= 0 {
		serviceName = fmt.Sprintf("ex-%v-%v-%v", kind, resourceName, ns)
	}
	args = append(args, "-n")
	args = append(args, ns)
	args = append(args, fmt.Sprintf(name, serviceName))
	args = append(args, fmt.Sprintf(targetPort, internalPort))
	args = append(args, fmt.Sprintf(externalPortTpl, externalPort))
	args = append(args, fmt.Sprintf(externalIp, nodeIp))
	args = append(args, fmt.Sprintf(typeArg, ExposeTypeNodePort))
	return KubeApi(cluster, args...)
}

// 服务暴露
type ExposeSvc struct {
	Name      string
	Namespace string
	Ip        string
	Port      string
	Selector  string
}

var exposeReg = regexp.MustCompile("^(\\d+)")

// 获取所有对外暴露的服务
func GetExpose(cluster KubeCluster, ns string) []ExposeSvc {

	var args []string
	args = append(args, CmdGet)
	args = append(args, "svc")
	if len(ns) <= 0 {
		ns = "default"
	}
	args = append(args, "-n")
	args = append(args, ns)
	args = append(args, "-o")
	args = append(args, "wide")
	res, err := KubeApi(cluster, args...)
	if err != nil {
		return nil
	}
	tableData := middleware.RenderTable(res)
	if tableData == nil || len(tableData) <= 0 {
		return nil
	}
	var result []ExposeSvc
	for _, row := range tableData {
		svcType := row["TYPE"]
		if svcType != "NodePort" {
			continue
		}
		externalIp := row["EXTERNAL-IP"]
		if len(externalIp) <= 0 || externalIp == "<none>" {
			continue
		}
		resRow := ExposeSvc{
			Ip:        externalIp,
			Name:      row["NAME"],
			Namespace: ns,
			Selector:  row["SELECTOR"],
		}
		ports := row["PORT(S)"] // 80:31099/TCP
		portRes := exposeReg.FindStringSubmatch(ports)
		if len(portRes) > 1 {
			resRow.Port = portRes[1]
		}
		result = append(result, resRow)
	}
	return result
}
