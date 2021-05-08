package proxy

import (
	"bulk/pkg/xcmd"

	"github.com/spf13/cobra"
)

var flags struct {
	Connect string   `shorthand:"c" usage:""`
	Listen  []string `shorthand:"l" usage:""`
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Listen on ports and proxy data to another connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if err := xcmd.Bind(cmd, &flags); err != nil {
		panic(err)
	}

	return cmd
}
