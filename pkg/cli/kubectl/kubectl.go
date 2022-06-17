package kubectl

import (
	"flag"
	"fmt"
	"k8s.io/kubectl/pkg/cmd/plugin"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/kubectl/pkg/cmd"
)

// Main kubectl main function.
// Borrowed from https://github.com/kubernetes/kubernetes/blob/master/cmd/kubectl/kubectl.go.
func Main() {
	rand.Seed(time.Now().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := EmbedCommand().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// EmbedCommand Used to embed the kubectl command.
func EmbedCommand() *cobra.Command {
	c := cmd.NewDefaultKubectlCommand()
	c.Short = "Kubectl controls the Kubernetes cluster manager"

	return c
}

// call comamnd kubectl drain nodename to drain all pods of nodename firstly then call kubectl delete node nodename to delete the node from cluster
//added by edgego
func KubectlDeleteNode(cluster, name string) error {
	contextArgs := []string{"config", "use-context", cluster}
	c := cmd.NewDefaultKubectlCommandWithArgs(cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes), contextArgs, os.Stdin, os.Stdout, os.Stderr)
	c.Short = "Kubectl controls the Kubernetes cluster manager"
	if err := c.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	drainArgs := []string{"drain", name, "--ignore-daemonsets", "--delete-local-data"}
	c = cmd.NewDefaultKubectlCommandWithArgs(cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes), drainArgs, os.Stdin, os.Stdout, os.Stderr)
	if err := c.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	deleteNodeArgs := []string{"delete", "node", name}
	c = cmd.NewDefaultKubectlCommandWithArgs(cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes), deleteNodeArgs, os.Stdin, os.Stdout, os.Stderr)
	if err := c.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	return nil
}
