//go:build plan9

package localcommand

import (
	"os"
	"strconv"
	"syscall"
)

func closeSignalFromInt(signal int) os.Signal {
	switch signal {
	case 2:
		return os.Interrupt
	case 9:
		return os.Kill
	default:
		return syscall.Note(strconv.Itoa(signal))
	}
}
