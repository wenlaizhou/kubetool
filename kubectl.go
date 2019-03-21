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

/*
kubectl controls the Kubernetes cluster manager.

Find more information at: https://kubernetes.io/docs/reference/kubectl/overview/

Basic Commands (Beginner):
  create         Create a resource from a file or from stdin.
  expose         使用 replication controller, service, deployment 或者 pod 并暴露它作为一个 新的
Kubernetes Service
  run            在集群中运行一个指定的镜像
  set            为 objects 设置一个指定的特征

Basic Commands (Intermediate):
  explain        查看资源的文档
  get            显示一个或更多 resources
  edit           在服务器上编辑一个资源
  delete         Delete resources by filenames, stdin, resources and names, or by resources and label selector

Deploy Commands:
  rollout        Manage the rollout of a resource
  scale          为 Deployment, ReplicaSet, Replication Controller 或者 Job 设置一个新的副本数量
  autoscale      自动调整一个 Deployment, ReplicaSet, 或者 ReplicationController 的副本数量

Cluster Management Commands:
  certificate    修改 certificate 资源.
  cluster-info   显示集群信息
  top            Display Resource (CPU/Memory/Storage) usage.
  cordon         标记 node 为 unschedulable
  uncordon       标记 node 为 schedulable
  drain          Drain node in preparation for maintenance
  taint          更新一个或者多个 node 上的 taints

Troubleshooting and Debugging Commands:
  describe       显示一个指定 resource 或者 group 的 resources 详情
  logs           输出容器在 pod 中的日志
  attach         Attach 到一个运行中的 container
  exec           在一个 container 中执行一个命令
  port-forward   Forward one or more local ports to a pod
  proxy          运行一个 proxy 到 Kubernetes API server
  cp             复制 files 和 directories 到 containers 和从容器中复制 files 和 directories.
  auth           Inspect authorization

Advanced Commands:
  diff           Diff live version against would-be applied version
  apply          通过文件名或标准输入流(stdin)对资源进行配置
  patch          使用 strategic merge patch 更新一个资源的 field(s)
  replace        通过 filename 或者 stdin替换一个资源
  wait           Experimental: Wait for a specific condition on one or many resources.
  convert        在不同的 API versions 转换配置文件

Settings Commands:
  label          更新在这个资源上的 labels
  annotate       更新一个资源的注解
  completion     Output shell completion code for the specified shell (bash or zsh)

Other Commands:
  api-resources  Print the supported API resources on the server
  api-versions   Print the supported API versions on the server, in the form of "group/version"
  config         修改 kubeconfig 文件
  plugin         Provides utilities for interacting with plugins.
  version        输出 client 和 server 的版本信息

Usage:
  kubectl [flags] [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
*/

/*
kubectl options
The following options can be passed to any command:

	 --alsologtostderr=false: log to standard error as well as files
	 --as='': Username to impersonate for the operation
	 --as-group=[]: Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
	 --cache-dir='/root/.kube/http-cache': Default HTTP cache directory
	 --certificate-authority='': Path to a cert file for the certificate authority
	 --client-certificate='': Path to a client certificate file for TLS
	 --client-key='': Path to a client key file for TLS
	 --cluster='': The name of the kubeconfig cluster to use
	 --context='': The name of the kubeconfig context to use
	 --insecure-skip-tls-verify=false: If true, the server's certificate will not be checked for validity. This will
make your HTTPS connections insecure
	 --kubeconfig='': Path to the kubeconfig file to use for CLI requests.
	 --log-backtrace-at=:0: when logging hits line file:N, emit a stack trace
	 --log-dir='': If non-empty, write log files in this directory
	 --log-file='': If non-empty, use this log file
	 --log-flush-frequency=5s: Maximum number of seconds between log flushes
	 --logtostderr=true: log to standard error instead of files
	 --match-server-version=false: Require server version to match client version
 -n, --namespace='': If present, the namespace scope for this CLI request
	 --profile='none': Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex)
	 --profile-output='profile.pprof': Name of the file to write the profile to
	 --request-timeout='0': The length of time to wait before giving up on a single server request. Non-zero values
should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.
 -s, --server='': The address and port of the Kubernetes API server
	 --skip-headers=false: If true, avoid header prefixes in the log messages
	 --stderrthreshold=2: logs at or above this threshold go to stderr
	 --token='': Bearer token for authentication to the API server
	 --user='': The name of the kubeconfig user to use
 -v, --v=0: log level for V logs
	 --vmodule=: comma-separated list of pattern=N settings for file-filtered logging
*/

// kubeproxy
// kubectl proxy --address='0.0.0.0' --port=8080 --www='' --www-prefix='/static' --accept-paths='^.*' --accept-hosts='.*'
// 为不同集群创建不同文件夹, 每个文件夹内置一个config , 自动生成 cache及http-cache文件夹进行缓存.
// 使用--kubeconfig 来指定不同集群的配置链接方式

// kcloud自身保存配置文件到对应config目录之中, 重启不丢失集群配置信息

// 祛除节点: k drain k8s-harbor --force --ignore-daemonsets --delete-local-data
// k delete node k8s-harbor

// pod 迁移到新node节点上

func init() {
}
