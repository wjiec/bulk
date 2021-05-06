package main

import (
	"bulk/internnal/bulk/crlf"
	"bulk/pkg/xcontext"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{Use: "bulk", Short: "A bulk tools", SilenceUsage: true, SilenceErrors: true}

	root.AddCommand(crlf.Command())

	if err := root.ExecuteContext(xcontext.WithExitSignal(context.Background())); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %s\n\ntry 'bulk --help' for more information.", err)
	}
}
