//go:build plan9

package localcommand

import (
	"io"
	"os"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var DefaultCloseSignal os.Signal = os.Interrupt

const (
	DefaultCloseTimeout = 10 * time.Second
)

type LocalCommand struct {
	command string
	argv    []string

	closeSignal  os.Signal
	closeTimeout time.Duration
}

func New(command string, argv []string, options ...Option) (*LocalCommand, error) {
	lcmd := &LocalCommand{
		command: command,
		argv:    argv,

		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,
	}

	for _, option := range options {
		option(lcmd)
	}

	return nil, unsupportedError()
}

func (lcmd *LocalCommand) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func (lcmd *LocalCommand) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func (lcmd *LocalCommand) Close() error {
	return unsupportedError()
}

func (lcmd *LocalCommand) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{
		"command": lcmd.command,
		"argv":    lcmd.argv,
		"pid":     0,
	}
}

func (lcmd *LocalCommand) ResizeTerminal(width int, height int) error {
	return unsupportedError()
}

func unsupportedError() error {
	return errors.Errorf("local command backend is unsupported on %s", runtime.GOOS)
}
