package kubetool

import (
	"encoding/json"
	"github.com/wenlaizhou/kubetype"
)

// 获取时间列表
func GetEvents(cluster KubeCluster) (kubetype.EventList, error) {
	resObj := kubetype.EventList{}
	res, err := KubeApi(cluster, CmdGet, "events", ArgsAllNamespaces, ArgsOutput, ArgsJson)
	if err != nil {
		return resObj, err
	}
	err = json.Unmarshal([]byte(res), &resObj)
	return resObj, err
}
