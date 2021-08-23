package kcobra

import (
	"github.com/spf13/cobra"
)

type Command struct {
	cobra.Command

	Current *cobra.Command
}

func NewRootCommand(name string) *Command {
	root := &Command{Command: cobra.Command{Use: name, SilenceUsage: true, SilenceErrors: true}}
	root.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		root.Current = cmd
	}

	return root
}
