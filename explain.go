package kubetool

// explain
func ExplainResource(resourceName string) (string, error) {
	for _, v := range Cluster {
		res, err := KubeApi(v, CmdExplain, resourceName, ArgRecursive)
		return res, err
	}
	return "", nil
}
