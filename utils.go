package kubetool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"io"
	"os/exec"
	"strings"
	"time"
)

var K8sLogger = middleware.GetLogger("k8s")

type KubectlResult struct {
	Header []string
	Table  []map[string]string
}

// 进行kubernetes接口调用
func ExecKubectl(cluster KubeCluster, commands ...string) (string, error) {
	var err error
	cachePath := fmt.Sprintf("--cache-dir=%s", cluster.CachePath)
	kubeConfig := fmt.Sprintf("--kubeconfig=%s", cluster.ConfPath)
	var cmd *exec.Cmd
	var args []string
	args = append(args, cachePath, kubeConfig)
	for _, c := range commands {
		args = append(args, c)
	}

	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)
	cmd = exec.CommandContext(ctx, Kubectl, args...)
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		K8sLogger.ErrorF("kubectl error : %s", err.Error())
		K8sLogger.ErrorF("%s, %v", cmd.Path, cmd.Args)
		return "", err
	}
	_, _ = io.Copy(&stdoutBuf, stdoutIn)
	_, _ = io.Copy(&stderrBuf, stderrIn)
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	if err != nil {
		K8sLogger.Error(err.Error())
	}
	_ = cmd.Wait()
	return fmt.Sprintf("%s\n%s", outStr, errStr), err
	// out, err := cmd.CombinedOutput()
	// return string(out), err
}

// 接口结构化调用
func KubectlStructExec(cmds ...string) (*KubectlResult, error) {
	output, err := exec.Command(cmds[0], cmds[1:]...).CombinedOutput()
	fmt.Printf("%#v", string(output))
	if err != nil {
		return nil, err
	}
	outputStr := string(output)
	outputLines := strings.Split(outputStr, "\n")
	// first line is data table
	if len(outputLines) <= 0 {
		return nil, errors.New("no result")
	}
	res := &KubectlResult{
		Header: []string{},
		Table:  []map[string]string{},
	}
	for _, h := range strings.Split(outputLines[0], " ") {
		if len(h) > 0 {
			res.Header = append(res.Header, h)
		}
	}
	for _, contentLines := range outputLines[1:] {
		table := make(map[string]string)
		headerIndex := 0
		for _, content := range strings.Split(contentLines, " ") {
			if len(content) > 0 {
				table[res.Header[headerIndex]] = content
				headerIndex++
			}
		}
		res.Table = append(res.Table, table)
	}
	return res, nil
}
