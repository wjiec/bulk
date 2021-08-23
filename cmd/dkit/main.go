package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wjiec/dkit/internal/version"

	"github.com/wjiec/dkit/internal/crlf"
	"github.com/wjiec/dkit/pkg/kcobra"
)

var (
	Name        = "dkit"
	Version     = "v0.0.0-dev"
	GitRevision = "0000000"
	BuildTime   = "2006/01/02"
)

func main() {
	root := kcobra.NewRootCommand("dkit")

	root.AddCommand(crlf.Command())
	root.AddCommand(version.Command(Name, Version, BuildTime, GitRevision))

	if err := root.ExecuteContext(context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %s\n\n", err)
		if root.Current == nil {
			_, _ = fmt.Fprint(os.Stderr, "try 'dkit --help' for more information.")
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "try 'dkit %s --help' for more information.", root.Current.Name())
		}
	}
}
