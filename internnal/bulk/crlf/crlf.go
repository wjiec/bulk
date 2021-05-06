package crlf

import (
	"bulk/pkg/xcmd"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var flags struct {
	DryRun    bool     `default:"false" usage:"just print a list of files that will be converted"`
	Target    string   `default:"unix" usage:"target newline style, optional values are unix and dos"`
	Recursive bool     `default:"false" shorthand:"r" usage:"recursively convert all files in the directories"`
	Exclude   []string `usage:"exclude files matching ·patterns·"`
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crlf [flags] files...",
		Short: "Convert text files between dos and unix newline style",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags.Target = strings.ToLower(flags.Target)
			if flags.Target != "unix" && flags.Target != "dos" {
				return errors.Errorf("unsupported newline style %s", flags.Target)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if err := xcmd.Bind(cmd, &flags); err != nil {
		panic(err)
	}

	return cmd
}

func convert(filename string) (bool, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	if flags.DryRun {
		return bytes.Contains(body, []byte("")), nil
	}
	return false, nil
}
