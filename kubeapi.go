package kubetool

import (
	"fmt"
	"github.com/wenlaizhou/middleware"
	"strings"
)

// 进行kubernetes接口调用
func KubeApi(cluster KubeCluster, commands ...string) (string, error) {
	cachePath := fmt.Sprintf("--cache-dir=%s", cluster.CachePath)
	kubeConfig := fmt.Sprintf("--kubeconfig=%s", cluster.ConfPath)
	var args []string
	args = append(args, cachePath, kubeConfig)
	for _, c := range commands {
		args = append(args, c)
	}
	K8sLogger.InfoF("k8s调用: %v", args)
	res, err := middleware.ExecCmdWithTimeout(5, Kubectl, args...)
	if err != nil {
		K8sLogger.Error(err.Error())
	}
	return strings.TrimSpace(res), err
}
