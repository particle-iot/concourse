//go:build windows
// +build windows

package pty

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}

func OpenRawTerm() (Term, error) {
	return noopRestoreTerm{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, nil
}

type noopRestoreTerm struct {
	io.Reader
	io.Writer
}

func (noopRestoreTerm) Restore() error { return nil }
