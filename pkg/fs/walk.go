package fs

import (
	"os"
	"path/filepath"
)

// VisitFiles visits the file tree rooted at root, calling fn for each file only
func VisitFiles(root string, fn func(filename string, stat os.FileInfo) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return fn(path, info)
		}
		return nil
	})
}
