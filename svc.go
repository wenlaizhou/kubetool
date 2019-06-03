package kubetool

import (
	"encoding/json"
	"fmt"
	"github.com/wenlaizhou/kubetype"
)

// 根据请求参数获取服务列表
//
// 改为多集群数据
func GetService(clusterName string, name string, namespace string, selector string, fieldSelector string) kubetype.ServiceList {
	res := kubetype.ServiceList{}
	cluster, hasCluster := Cluster[clusterName]
	if !hasCluster {
		return res
	}
	var args []string
	args = append(args, CmdGet)
	args = append(args, "svc")

	single := len(name) > 0

	if single {
		args = append(args, name)
		if len(namespace) <= 0 { // 单服务需要指定命名空间
			return res
		}
		args = append(args, "-n", namespace)
	} else {
		if len(namespace) <= 0 {
			args = append(args, "--all-namespaces")
		} else {
			args = append(args, "-n", namespace)
		}
	}

	if len(selector) > 0 {
		args = append(args, fmt.Sprintf("--selector='%s'", selector))
	}
	if len(fieldSelector) > 0 {
		args = append(args, fmt.Sprintf("--field-selector='%s'", selector))
	}

	args = append(args, "-o")
	args = append(args, "json")

	cmdRes, err := KubeApi(cluster, args...)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get svc error : %s, name: %s, namespace: %s, selector: %s, fieldSelector: %s",
			cluster.Name, err.Error(), name, namespace, selector, fieldSelector)
		return res
	}
	if single {
		svc := kubetype.Service{}
		err = json.Unmarshal([]byte(cmdRes), &svc)
		res.Items = append(res.Items, svc)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get svc error json parse error : %s, name: %s, namespace: %s, selector: %s, fieldSelector: %s",
				cluster.Name, err.Error(), name, namespace, selector, fieldSelector)
			return res
		}
	} else {
		err = json.Unmarshal([]byte(cmdRes), &res)
		if err != nil {
			K8sLogger.ErrorF("cluster: %s get svc error json parse error : %s, name: %s, namespace: %s, selector: %s, fieldSelector: %s",
				cluster.Name, err.Error(), name, namespace, selector, fieldSelector)
			return res
		}
	}
	return res
}
