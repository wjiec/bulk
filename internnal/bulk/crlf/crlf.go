package crlf

import (
	"bulk/pkg/fs"
	"bulk/pkg/utils"
	"bulk/pkg/xcmd"
	"bulk/pkg/xsync"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var flags struct {
	DryRun    bool   `default:"false" usage:"just print a list of files that will be converted"`
	Target    string `default:"unix" usage:"target newline style, optional values are unix and dos"`
	Recursive bool   `default:"false" shorthand:"r" usage:"recursively convert all files in the directories"`
}

var (
	// UNIX newline style by default
	from, to = []byte("\r\n"), []byte("\n")
	// ErrNotRegularTextFile represents file is not a regular text file
	ErrNotRegularTextFile = errors.New("not regular text file")
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crlf [flags] files...",
		Short: "Convert text files between dos and unix newline style",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags.Target = strings.ToLower(flags.Target)
			if flags.Target != "unix" && flags.Target != "dos" {
				return errors.Errorf("unsupported newline style %s", flags.Target)
			} else if flags.Target == "dos" {
				from, to = to, from
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var wg xsync.WaitGroup
			for tuple := range walk(args) {
				wg.Do(func(args ...interface{}) {
					if args[1] != nil {
						fmt.Printf("[ERROR]\t%s\t:\t%s\n", args[0], args[1])
						return
					}

					if ok, err := convert(args[0].(string)); err != nil {
						if err == ErrNotRegularTextFile && flags.DryRun {
							fmt.Printf("[SKIPPED]\t%s\t:\t%s\n", args[0], err)
						} else if err != ErrNotRegularTextFile {
							fmt.Printf("[ERROR]\t%s\t:\t%s\n", args[0], err)
						}
					} else if !ok && flags.DryRun {
						fmt.Printf("[SKIPPED]\t%s\t:\tdon't need to convert\n", args[0])
					} else if flags.DryRun {
						fmt.Printf("[SUCCESS]\t%s\n", args[0])
					}
				}, tuple.First, tuple.Second)
			}
			wg.Wait()
			return nil
		},
	}

	if err := xcmd.Bind(cmd, &flags); err != nil {
		panic(err)
	}

	return cmd
}

func walk(paths []string) <-chan *utils.Tuple {
	ch := make(chan *utils.Tuple)
	go func() {
		for _, p := range paths {
			var ok bool
			var err error

			if ok, err = fs.IsDirectory(p); ok {
				err = fs.VisitFiles(p, func(filename string, stat os.FileInfo) error {
					ch <- utils.NewTuple(filename, nil)
					return nil
				})
			} else if ok, err = fs.IsRegularFile(p); ok {
				ch <- utils.NewTuple(p, nil)
			}

			if err != nil {
				ch <- utils.NewTuple(p, err)
			}
		}
		close(ch)
	}()
	return ch
}

// convert convert file between unix and dos newline style
func convert(filename string) (bool, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, errors.Wrap(err, "unreadable file")
	}

	// not regular text file
	if bytes.Contains(bs, []byte{0}) {
		return false, ErrNotRegularTextFile
	}
	// don't need to convert
	if !bytes.Contains(bs, from) {
		return false, nil
	}

	if flags.DryRun {
		return true, nil
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return false, errors.Wrap(err, "unable to access file permission")
	}

	return true, ioutil.WriteFile(filename, bytes.ReplaceAll(bs, from, to), stat.Mode())
}
