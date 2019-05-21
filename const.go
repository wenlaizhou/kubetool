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

const ArgRecursive = "--recursive=true"

const ArgsOverwrite = "--overwrite"

const ArgsAllContainers = "--all-containers=true"

const ArgsNamespace = "-n"
