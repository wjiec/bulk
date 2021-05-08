package main

import (
	"bulk/internnal/bulk/crlf"
	"bulk/internnal/bulk/proxy"
	"bulk/pkg/xcontext"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var command *cobra.Command

func main() {
	root := cobra.Command{Use: "bulk", Short: "A bulk tools", SilenceUsage: true, SilenceErrors: true}

	root.AddCommand(crlf.Command())
	root.AddCommand(proxy.Command())
	root.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		command = cmd
	}

	if err := root.ExecuteContext(xcontext.WithExitSignal(context.Background())); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %s\n\ntry 'bulk %s --help' for more information.", err, command.Name())
	}
}
