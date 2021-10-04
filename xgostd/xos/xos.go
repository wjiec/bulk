package xos

import "os"

const (
	// ModeDir represents a file mode in *nix is rwx-rx-rx
	ModeDir = os.FileMode(0755)
	// ModeRegularFile represents a file mode in *nix is rw-r-r
	ModeRegularFile = os.FileMode(0644)
)
