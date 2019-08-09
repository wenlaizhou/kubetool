package kubetool

import (
	"errors"
	"fmt"
	"github.com/wenlaizhou/yml"
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

const confDir = "kubeconf"

type KubeYml struct {
	ApiVersion string `json:"apiVersion"`
}

// 初始化配置信息
func init() {
	currentDir, _ = os.Getwd()
	if !middleware.Exists(confDir) {
		middleware.Mkdir(confDir)
		return
	}
	Init("conf")
}

// 当前目录进行配置扫描
func Init(confPath string) {
	_ = filepath.Walk(confPath, walkPath)
}

// 加入集群
// 配置 或者 配置文件路径
// 集群名称
func JoinCluster(confPath string) error {
	clusterName := strings.Replace(confPath, ".config", "", -1)
	// 重复集群
	cachePath := fmt.Sprintf("%s/%s/%s/cache", currentDir, confPath, clusterName)
	// hasCluster, _ := Cluster[clusterName]
	conf := middleware.ReadString(confPath)
	if len(conf) <= 0 {
		return errors.New("配置为空")
	}
	Cluster[clusterName] = KubeCluster{
		Name:      clusterName,
		ConfPath:  confPath,
		CachePath: cachePath,
		Conf:      conf,
	}
	middleware.Mkdir(cachePath)
	return nil
}

// 动态新增集群
// 集群名称, 集群配置文件
func NewCluster(name string, conf string) error {
	if len(name) <= 0 || len(conf) <= 0 {
		return errors.New("配置数据为空")
	}
	confStruct := KubeYml{}
	err := yml.Unmarshal([]byte(conf), &confStruct)
	if err != nil {
		return err
	}
	confPath := fmt.Sprintf("%s/%s/%s.config", currentDir, confDir, name)
	cachePath := fmt.Sprintf("%s/%s/%s/cache", currentDir, confDir, name)
	_, _ = middleware.WriteString(
		fmt.Sprintf("%s/%s/%s.config", currentDir, confPath, name), conf)
	Cluster[name] = KubeCluster{
		Name:      name,
		ConfPath:  confPath,
		CachePath: cachePath,
		Conf:      conf,
	}
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
