package kubetool

// explain
func ExplainResource(resourceName string) (string, error) {
	for _, v := range Cluster {
		res, err := ExecKubectl(v, CmdExplain, resourceName, ArgRecursive)
		return res, err
	}
	return "", nil
}
