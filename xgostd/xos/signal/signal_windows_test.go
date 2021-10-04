// +build windows

package signal

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

var (
	sigUsr1 = syscall.SIGINT
	sigUsr2 = syscall.SIGINT
)

// sendSignal send a signal to the current process
// Note: windows platform may not support sending signals to process
func sendSignal(_ os.Signal) error {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}

	proc, err := dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}

	ret, _, err := proc.Call(syscall.CTRL_BREAK_EVENT, uintptr(os.Getpid()))
	if ret == 0 {
		return errors.Wrap(err, "GenerateConsoleCtrlEvent")
	}
	return nil
}
