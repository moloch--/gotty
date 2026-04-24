//go:build !plan9

package localcommand

import (
	"os"
	"syscall"
)

func closeSignalFromInt(signal int) os.Signal {
	return syscall.Signal(signal)
}
