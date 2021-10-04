// +build linux darwin openbsd

package signal

import (
	"os"
	"syscall"
)

var (
	sigUsr1 = syscall.SIGUSR1
	sigUsr2 = syscall.SIGUSR2
)

// sendSignal send a signal to the current process
func sendSignal(signal os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	return p.Signal(signal)
}
