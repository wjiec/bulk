package crlf

import (
	"bulk/pkg/fs"
	"bulk/pkg/xcmd"
	"bulk/pkg/xsync"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

// UNIX newline style by default
var from, to = []byte("\r\n"), []byte("\n")

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
			for filename := range walk(args) {
				filename := filename
				wg.Do(func() {
					if ok, err := convert(filename); err != nil {
						fmt.Printf("%s(error): %s\n", filename, err)
					} else if ok {
						fmt.Printf("%s: ok\n", filename)
					} else {
						fmt.Printf("%s(skipped): not regular text file\n", filename)
					}
				})
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

func walk(paths []string) <-chan string {
	ch := make(chan string)
	go func() {
		for _, path := range paths {
			if isDir, err := fs.IsDirectory(path); isDir {
				if flags.Recursive {
					_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						if info.Mode().IsRegular() {
							ch <- path
						}
						return nil
					})
				} else {
					fmt.Printf("%s(skipped): directory, use '-r' to recursively convert\n", path)
				}
			} else if err != nil {
				fmt.Printf("%s(error): %s\n", path, err)
			} else if isFile, err := fs.IsRegularFile(path); isFile {
				ch <- path
			} else {
				fmt.Printf("%s(error): %s\n", path, err)
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

	// not text file or don't need to convert
	if bytes.Contains(bs, []byte{0}) || !bytes.Contains(bs, from) {
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
