package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/wjiec/dkit/internal/crlf"
	"github.com/wjiec/dkit/internal/version"
	"github.com/wjiec/dkit/pkg/kcobra"
)

var (
	Name        = "dkit"
	Version     = "v0.0.0-dev"
	GitRevision = "0000000"
	BuildTime   = "2006/01/02"

	exitSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
)

func main() {
	root := &cobra.Command{
		Use:           Name,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(crlf.Command())
	root.AddCommand(version.Command(Name, Version, BuildTime, GitRevision))
	if err := kcobra.ExecuteWithSignal(root, exitSignals...); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
