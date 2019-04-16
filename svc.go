package kubetool

import (
	"encoding/json"
	"fmt"
	"github.com/wenlaizhou/kubetype"
)

// 根据请求参数获取服务列表
func GetService(cluster KubeCluster, name string, namespace string, selector string, fieldSelector string) kubetype.ServiceList {
	res := kubetype.ServiceList{}
	var args []string
	args = append(args, CmdGet)
	args = append(args, "svc")
	if len(name) > 0 {
		args = append(args, name)
	}
	args = append(args, "-o")
	args = append(args, "json")
	if len(namespace) <= 0 {
		args = append(args, "--all-namespaces")
	} else {
		args = append(args, "-n", namespace)
	}
	if len(selector) > 0 {
		args = append(args, fmt.Sprintf("--selector='%s'", selector))
	}
	if len(fieldSelector) > 0 {
		args = append(args, fmt.Sprintf("--field-selector='%s'", selector))
	}
	cmdRes, err := ExecKubectl(cluster, args...)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get svc error : %s, name: %s, namespace: %s, selector: %s, fieldSelector: %s",
			cluster.Name, err.Error(), name, namespace, selector, fieldSelector)
		return res
	}
	err = json.Unmarshal([]byte(cmdRes), &res)
	if err != nil {
		K8sLogger.ErrorF("cluster: %s get svc error : %s, name: %s, namespace: %s, selector: %s, fieldSelector: %s",
			cluster.Name, err.Error(), name, namespace, selector, fieldSelector)
		return res
	}
	return res
}
