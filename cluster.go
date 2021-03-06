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

var confDir = "kubeconf"

type KubeYml struct {
	ApiVersion string `yml:"apiVersion"`
	Kind       string `yml:"kind"`
}

// 初始化配置信息
func init() {
	currentDir, _ := os.Getwd()
	confDir = fmt.Sprintf("%v/%v", currentDir, confDir)
	if !middleware.Exists(confDir) {
		middleware.Mkdir(confDir)
		return
	}
	Init()
}

// 当前目录进行配置扫描
func Init() {
	_ = filepath.Walk(confDir, walkPath)
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
	confPath := fmt.Sprintf("%s/%s.config", confDir, name)
	cachePath := fmt.Sprintf("%s/%s/cache", confDir, name)
	K8sLogger.InfoF("创建新集群: %v", name)
	_, _ = middleware.WriteString(confPath, conf)
	Cluster[name] = KubeCluster{
		Name:      name,
		ConfPath:  confPath,
		CachePath: cachePath,
		Conf:      conf,
	}
	return nil
}

// 漫游配置路径
func walkPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info == nil {
		return errors.New("file not exist")
	}
	if !strings.HasSuffix(info.Name(), ".config") {
		return nil
	}
	clusterName := strings.Replace(info.Name(), ".config", "", -1)
	cachePath := fmt.Sprintf("%s/%s/cache", confDir, clusterName)
	confPath := fmt.Sprintf("%s/%s", confDir, info.Name())
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
