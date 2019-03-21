package kubetool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wenlaizhou/middleware"
)

// 集群配置
var Cluster = make(map[string]KubeCluster)

// 集群配置结构
type KubeCluster struct {
	Name      string `json:"name"`
	ConfPath  string `json:"confPath"`
	CachePath string `json:"cachePath"`
	Conf      string `json:"conf"`
}

var currentDir string

// 初始化配置信息
func init() {

	currentDir, _ = os.Getwd()
	if !middleware.Exists("conf") {
		return
	}
	_ = filepath.Walk("conf", walkPath)
}

// 漫游配置路径
func walkPath(path string, info os.FileInfo, err error) error {
	if !strings.HasSuffix(info.Name(), ".config") {
		return nil
	}
	clusterName := strings.Replace(info.Name(), ".config", "", -1)
	cachePath := fmt.Sprintf("%s/conf/%s/cache", currentDir, clusterName)
	confPath := fmt.Sprintf("%s/conf/%s", currentDir, info.Name())
	conf := middleware.ReadString(path)
	Cluster[clusterName] = KubeCluster{
		Name:      clusterName,
		ConfPath:  confPath,
		CachePath: cachePath,
		Conf:      conf,
	}
	middleware.Mkdir(cachePath)
	return nil
}
