package kubetool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wenlaizhou/kubetype"
)

// 获取资源列表
func GetResourceList(clusterName string, resourceName string, ns string) (interface{}, error) {
	if len(resourceName) <= 0 {
		return nil, errors.New("resourceName 为空")
	}
	if len(clusterName) <= 0 {
		return nil, errors.New("集群参数错误 为空")
	}
	cluster, hasData := Cluster[clusterName]
	if !hasData {
		return nil, errors.New(fmt.Sprintf("集群参数错误 没有该集群%v", clusterName))
	}
	args := []string{
		CmdGet,
		resourceName,
	}
	if len(ns) > 0 {
		args = append(args, ArgsNamespace, ns)
	} else {
		args = append(args, ArgsAllNamespaces)
	}
	args = append(args, ArgsOutput, ArgsJson)
	res, err := KubeApi(cluster, args...)
	if err != nil {
		return nil, err
	}
	switch resourceName {
	case "po":
		resObj := kubetype.PodList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "cm":
		resObj := kubetype.ConfigMapList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "ev":
		resObj := kubetype.EventList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "ep":
		resObj := kubetype.EndpointsList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "ns":
		resObj := kubetype.NamespaceList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "pv":
		resObj := kubetype.PersistentVolumeList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "pvc":
		resObj := kubetype.PersistentVolumeClaimList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "svc":
		resObj := kubetype.ServiceList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "ds":
		resObj := kubetype.DaemonSetList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "rs":
		resObj := kubetype.ReplicaSetList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "sts":
		resObj := kubetype.StatefulSetList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "sc":
		resObj := kubetype.StorageClassList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "deploy":
		resObj := kubetype.DeploymentList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "ing":
		resObj := kubetype.IngressList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "secrets":
		resObj := kubetype.SecretList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "sa":
		resObj := kubetype.ServiceAccountList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "clusterroles":
		resObj := kubetype.ClusterRoleList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "clusterrolebindings":
		resObj := kubetype.ClusterRoleBindingList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "roles":
		resObj := kubetype.RoleList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	case "rolebindings":
		resObj := kubetype.RoleBindingList{}
		err = json.Unmarshal([]byte(res), &resObj)
		return resObj, err
	}
	return nil, errors.New(fmt.Sprintf("没有这种资源: %s", resourceName))
}

func GetResoureYaml(cluster KubeCluster, resourceName string, name string, namespace string) (string, error) {
	return KubeApi(cluster, CmdGet, resourceName, name, ArgsNamespace, namespace, ArgsOutput, ArgsYml)
}

func DescResource(cluster KubeCluster, resourceName string, name string, namespace string) (string, error) {
	return KubeApi(cluster, CmdDesc, resourceName, name, ArgsNamespace, namespace)
}

func DeleteResource(cluster KubeCluster, resourceName string, name string, namespace string, force bool) error {

	args := []string{CmdDelete}
	args = append(args, resourceName, name, ArgsNamespace, name, "-R", "--wait=false")
	if force {
		args = append(args, "--grace-period=0", "--force")
	}
	// res, err := KubeApi(cluster, CmdDelete, resourceName, name, ArgsNamespace, namespace, "-R", "--wait=false")
	res, err := KubeApi(cluster, args...)

	if err == nil {
		K8sLogger.InfoF("%s: 删除k8s资源%s %s:%s, 结果为: %s", cluster.Name, resourceName, name, namespace, res)
	} else {
		K8sLogger.ErrorF("%s: 删除k8s资源%s %s:%s, 错误: %s, 结果为: %s", cluster.Name, resourceName, name, namespace, err.Error(), res)
	}
	return err
}

// 批量删除资源数据结构
type DeleteResStruct struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// 批量删除资源
func DeleteResourceList(cluster KubeCluster, resourceName string, resources []DeleteResStruct) {
	for _, res := range resources {
		resource := res
		resourceName := resourceName
		cluster := cluster
		go DeleteResource(cluster, resourceName, resource.Name, resource.Namespace)
	}
}
