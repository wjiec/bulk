package kcobra

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wjiec/dkit/pkg/signal"
)

type executeError struct {
	raw       error
	root, sub string
}

func (e *executeError) Error() string {
	return fmt.Sprintf("fatal: %s\n\ntry '%s %s --help' for more information.",
		e.raw, e.root, e.sub)
}

func injectCommand(cmd *cobra.Command, applier func(cmd *cobra.Command)) {
	cmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		applier(command)
		return err
	})

	backupPersistentPreRun := cmd.PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		applier(cmd)
		if backupPersistentPreRun != nil {
			backupPersistentPreRun(cmd, args)
		}
	}
}

func ExecuteWithSignal(cmd *cobra.Command, signals ...os.Signal) error {
	ctx, cancel := signal.WithContext(context.Background(), signals...)
	defer cancel()

	subCommand := "\b"
	injectCommand(cmd, func(cmd *cobra.Command) { subCommand = cmd.Name() })
	if err := cmd.ExecuteContext(ctx); err != nil {
		return &executeError{root: cmd.Name(), raw: err, sub: subCommand}
	}
	return nil
}
