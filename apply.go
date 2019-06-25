package kubetool

import (
	"fmt"
	"math/rand"
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

// 部署服务
func Apply(cluster KubeCluster, content string, update bool) (string, error) {
	fileName := fmt.Sprintf("%s/%s/%s_%d.yaml", CurrentDir, DeployYamlPath, time.Now().Format(TimeStr), rand.Int()%100)
	_, _ = middleware.WriteString(fileName, content)
	if update {
		return KubeApi(cluster, CmdApply, "-f", fileName, "--cascade=true") // , "--prune=true", "--all")
	} else {
		return KubeApi(cluster, CmdApply, "-f", fileName)
	}
}
