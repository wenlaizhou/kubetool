package kubetool

const Kubectl = "kubectl"

// yaml conf
const CmdCreate = "create"
const CmdApply = "apply"
const CmdReplace = "replace"

// service
const CmdRun = "run"
const CmdSet = "set"
const CmdExpose = "expose"

// res control
const CmdExplain = "explain"
const CmdGet = "get"
const CmdEdit = "edit"
const CmdDelete = "delete"

// scale
const CmdRollout = "rollout"
const CmdScale = "scale"
const CmdAutoscale = "autoscale"

// cluster
const CmdClusterInfo = "cluster-info"
const CmdTop = "top"

// node
const CmdCordon = "cordon"
const CmdUncordon = "uncordon"
const CmdDrain = "drain"
const Cmdtrait = "taint"

// control
const CmdDesc = "describe"
const CmdLogs = "logs"
const CmdExec = "exec"
const CmdProxy = "proxy"
const CmdCp = "cp" // 复制 files 和 directories 到 containers 和从容器中复制 files 和 directories.

// label
const CmdLabel = "label"         // 更新在这个资源上的 labels
const CmdAnnotate = "annotation" // 更新一个资源的注解

// api
const CmdResource = "api-resources"
const CmdApiVersion = "api-versions"
const CmdVersion = "version"

// args
const ArgRecursive = "--recursive=true"
const ArgsOverwrite = "--overwrite"
const ArgsAllContainers = "--all-containers=true"

// namespace
const ArgsNamespace = "-n"
const ArgsAllNamespaces = "--all-namespaces"

// output
const ArgsOutput = "-o"
const ArgsJson = "json"
const ArgsYml = "yaml"

var resourceTypes = []string{
	"pod", "service", "deployment", "ns", "ds", "sts", "cm",
}

// 判断资源类型是否是合法类型
func IsResourceType(rt string) bool {
	if len(rt) <= 0 {
		return false
	}
	for _, v := range resourceTypes {
		if v == rt {
			return true
		}
	}
	return false
}
