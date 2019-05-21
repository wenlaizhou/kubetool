package kubetool

import (
	"fmt"
	"os"
	"time"

	"github.com/wenlaizhou/middleware"
)

var CurrentDir string

// 时间格式
const TimeStr = "2006_1_2_15_04_05"

// 文件目录
const DeployYamlPath = "apply"

func init() {
	// 初始化文件保存目录
	CurrentDir, _ = os.Getwd()
	middleware.Mkdir(DeployYamlPath)
}

// 替换服务
func Replace(content string) (string, error) {
	timeStr := time.Now().Format("2006_1_2_15_04_05")
	middleware.Mkdir("deploy")
	_, _ = middleware.WriteString(fmt.Sprintf("deploy/%s.yaml", timeStr), content)
	return "", nil
}

// 部署服务
func Apply(cluster KubeCluster, content string, update bool) (string, error) {
	fileName := fmt.Sprintf("%s/%s/%s.yaml", CurrentDir, DeployYamlPath, time.Now().Format(TimeStr))
	_, _ = middleware.WriteString(fileName, content)
	if update {
		return KubeApi(cluster, CmdApply, "-f", fileName, "--cascade=true") // , "--prune=true", "--all")
	} else {
		return KubeApi(cluster, CmdApply, "-f", fileName)
	}
}
