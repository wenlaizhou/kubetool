package kubetool

import "fmt"

/*
Examples:
  # Scale a replicaset named 'foo' to 3.
  kubectl scale --replicas=3 rs/foo

  # Scale a resource identified by type and name specified in "foo.yaml" to 3.
  kubectl scale --replicas=3 -f foo.yaml

  # If the deployment named mysql's current size is 2, scale mysql to 3.
  kubectl scale --current-replicas=2 --replicas=3 deployment/mysql

  # Scale multiple replication controllers.
  kubectl scale --replicas=5 rc/foo rc/bar rc/baz

  # Scale statefulset named 'web' to 3.
  kubectl scale --replicas=3 statefulset/web

Options:
      --all=false: Select all resources in the namespace of the specified resource types
      --allow-missing-template-keys=true: If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats.
      --current-replicas=-1: Precondition for current size. Requires that the current size of the resource match this value in order to scale.
  -f, --filename=[]: Filename, directory, or URL to files identifying the resource to set a new size
  -o, --output='': Output format. One of: json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-file.
      --record=false: Record current kubectl command in the resource annotation. If set to false, do not record the command. If set to true, record the command. If not set, default to updating the existing annotation value only if one already exists.
  -R, --recursive=false: Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
      --replicas=0: The new desired number of replicas. Required.
      --resource-version='': Precondition for resource version. Requires that the current resource version match this value in order to scale.
  -l, --selector='': Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
      --template='': Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --timeout=0s: The length of time to wait before giving up on a scale operation, zero means don't wait. Any other values should contain a corresponding time unit (e.g. 1s, 2m, 3h).

Usage:
  kubectl scale [--resource-version=version] [--current-replicas=count] --replicas=COUNT (-f FILENAME | TYPE NAME) [options]
*/

func Scale(cluster KubeCluster, resourceName string, name string, scaleNumber int, ns string) (string, error) {
	var args []string
	args = append(args, CmdScale)
	args = append(args, fmt.Sprintf("--replicas=%v", scaleNumber))
	args = append(args, fmt.Sprintf("%v/%v", resourceName, name))
	return ExecKubectl(cluster, args...)
}

/*
Expose a resource as a new Kubernetes service.

Looks up a deployment, service, replica set, replication controller or pod by name and uses the selector for that
resource as the selector for a new service on the specified port. A deployment or replica set will be exposed as a
service only if its selector is convertible to a selector that service supports, i.e. when the selector contains only
the matchLabels component. Note that if no port is specified via --port and the exposed resource has multiple ports, all
will be re-used by the new service. Also if no labels are specified, the new service will re-use the labels from the
resource it exposes.

Possible resources include (case insensitive):

pod (po), service (svc), replicationcontroller (rc), deployment (deploy), replicaset (rs)

Examples:
  # Create a service for a replicated nginx, which serves on port 80 and connects to the containers on port 8000.
  kubectl expose rc nginx --port=80 --target-port=8000

  # Create a service for a replication controller identified by type and name specified in "nginx-controller.yaml",
which serves on port 80 and connects to the containers on port 8000.
  kubectl expose -f nginx-controller.yaml --port=80 --target-port=8000

  # Create a service for a pod valid-pod, which serves on port 444 with the name "frontend"
  kubectl expose pod valid-pod --port=444 --name=frontend

  # Create a second service based on the above service, exposing the container port 8443 as port 443 with the name
"nginx-https"
  kubectl expose service nginx --port=443 --target-port=8443 --name=nginx-https

  # Create a service for a replicated streaming application on port 4100 balancing UDP traffic and named 'video-stream'.
  kubectl expose rc streamer --port=4100 --protocol=udp --name=video-stream

  # Create a service for a replicated nginx using replica set, which serves on port 80 and connects to the containers on
port 8000.
  kubectl expose rs nginx --port=80 --target-port=8000

  # Create a service for an nginx deployment, which serves on port 80 and connects to the containers on port 8000.
  kubectl expose deployment nginx --port=80 --target-port=8000

Options:
      --allow-missing-template-keys=true: If true, ignore any errors in templates when a field or map key is missing in
the template. Only applies to golang and jsonpath output formats.
      --cluster-ip='': ClusterIP to be assigned to the service. Leave empty to auto-allocate, or set to 'None' to create
a headless service.
      --dry-run=false: If true, only print the object that would be sent, without sending it.
      --external-ip='': Additional external IP address (not managed by Kubernetes) to accept for the service. If this IP
is routed to a node, the service can be accessed by this IP in addition to its generated service IP.
  -f, --filename=[]: Filename, directory, or URL to files identifying the resource to expose a service
      --generator='service/v2': 使用 generator 的名称. 这里有 2 个 generators: 'service/v1' 和 'service/v2'.
为一个不同地方是服务端口在 v1 的情况下叫 'default', 如果在 v2 中没有指定名称.
默认的名称是 'service/v2'.
  -l, --labels='': Labels to apply to the service created by this call.
      --load-balancer-ip='': IP to assign to the LoadBalancer. If empty, an ephemeral IP will be created and used
(cloud-provider specific).
      --name='': 名称为最新创建的对象.
  -o, --output='': Output format. One of:
json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-file.
      --overrides='': An inline JSON override for the generated object. If this is non-empty, it is used to override the
generated object. Requires that the object supply a valid apiVersion field.
      --port='': 服务的端口应该被指定. 如果没有指定, 从被创建的资源中复制
      --protocol='': 创建 service 的时候伴随着一个网络协议被创建. 默认是 'TCP'.
      --record=false: Record current kubectl command in the resource annotation. If set to false, do not record the
command. If set to true, record the command. If not set, default to updating the existing annotation value only if one
already exists.
  -R, --recursive=false: Process the directory used in -f, --filename recursively. Useful when you want to manage
related manifests organized within the same directory.
      --save-config=false: If true, the configuration of current object will be saved in its annotation. Otherwise, the
annotation will be unchanged. This flag is useful when you want to perform kubectl apply on this object in the future.
      --selector='': A label selector to use for this service. Only equality-based selector requirements are supported.
If empty (the default) infer the selector from the replication controller or replica set.)
      --session-affinity='': If non-empty, set the session affinity for the service to this; legal values: 'None',
'ClientIP'
      --target-port='': Name or number for the port on the container that the service should direct traffic to.
Optional.
      --template='': Template string or path to template file to use when -o=go-template, -o=go-template-file. The
template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --type='': Type for this service: ClusterIP, NodePort, LoadBalancer, or ExternalName. Default is 'ClusterIP'.

Usage:
  kubectl expose (-f FILENAME | TYPE NAME) [--port=port] [--protocol=TCP|UDP|SCTP] [--target-port=number-or-name]
[--name=name] [--external-ip=external-ip-of-service] [--type=type] [options]

Use "kubectl options" for a list of global command-line options (applies to all commands).
*/
// 对外提供服务
// 资源类型: pod (po), service (svc),
// replicationcontroller (rc),
// deployment (deploy), replicaset (rs)

const ExposeTypeClusterIP = "ClusterIP" // default
const ExposeTypeNodePort = "NodePort"
const ExposeTypeLoadBalancer = "LoadBalancer"
const ExposeTypeExternalName = "ExternalName"
const typeArg = "--type=%s"
const externalPortTpl = "--port=%s"
const saveConfig = "--save-config=%s" // true | false A label selector to use for this service.
// Only equality-based selector requirements are supported. If empty (the default) infer the selector
// from the replication controller or replica set.)
const targetPort = "--target-port=%s"
const name = "--name=%s"
const externalIp = "--external-ip=%s"
const selectorTpl = "--selector=%s"

func Expose(cluster KubeCluster, kind string, resourceName string,
	serviceName string, ns string, externalPort string, existPort string,
	nodeIp string,
	selector string) (string, error) {
	// ExternalName - Exposes the Service using an arbitrary name
	// (specified by externalName in the spec) by returning a CNAME
	// record with the name. No proxy is used. This type requires v1.7 or higher of kube-dns.
	var args []string
	args = append(args, CmdExpose)
	args = append(args, kind)
	args = append(args, resourceName)
	if len(ns) <= 0 {
		ns = "default"
	}
	args = append(args, "-n")
	args = append(args, ns)
	args = append(args, fmt.Sprintf(name, serviceName))
	args = append(args, fmt.Sprintf(targetPort, existPort))
	args = append(args, fmt.Sprintf(externalPortTpl, externalPort))
	args = append(args, fmt.Sprintf(externalIp, nodeIp))
	args = append(args, fmt.Sprintf(typeArg, ExposeTypeNodePort))
	if len(selector) >= 0 {
		args = append(args, fmt.Sprintf(selectorTpl, selector))
	}
	return ExecKubectl(cluster, args...)
}
